package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type GitHubConfig struct {
	Username  string    `json:"username"`
	Token     string    `json:"token"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DBService struct {
	db            *sql.DB
	encryptionKey []byte
}

// NewDBService initializes PostgreSQL connection and tables
func NewDBService(databaseURL string) (*DBService, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// Load encryption key
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		keyStr = "default-32-byte-encryption-key!!"
		log.Println("WARNING: ENCRYPTION_KEY environment variable is not set. Using default development key.")
	}

	encryptionKey := []byte(keyStr)
	if len(encryptionKey) != 32 {
		db.Close()
		return nil, errors.New("ENCRYPTION_KEY must be exactly 32 bytes long")
	}

	svc := &DBService{
		db:            db,
		encryptionKey: encryptionKey,
	}

	// Run table migrations
	if err := svc.initTables(); err != nil {
		db.Close()
		return nil, err
	}

	return svc, nil
}

// Close closes the database connection
func (s *DBService) Close() error {
	return s.db.Close()
}

// initTables creates the github_configs table if it does not exist
func (s *DBService) initTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS github_configs (
		username VARCHAR(255) PRIMARY KEY,
		encrypted_token TEXT NOT NULL,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := s.db.Exec(query)
	return err
}

// SaveConfig encrypts the token and stores/updates it in the database
func (s *DBService) SaveConfig(username, token string) error {
	encryptedToken, err := Encrypt(token, s.encryptionKey)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO github_configs (username, encrypted_token, updated_at)
	VALUES ($1, $2, NOW())
	ON CONFLICT (username)
	DO UPDATE SET encrypted_token = EXCLUDED.encrypted_token, updated_at = NOW();`

	_, err = s.db.Exec(query, username, encryptedToken)
	return err
}

// GetConfig retrieves the configuration for a user and decrypts the token
func (s *DBService) GetConfig(username string) (*GitHubConfig, error) {
	query := `SELECT username, encrypted_token, updated_at FROM github_configs WHERE username = $1`
	row := s.db.QueryRow(query, username)

	var cfg GitHubConfig
	var encryptedToken string
	err := row.Scan(&cfg.Username, &encryptedToken, &cfg.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	decryptedToken, err := Decrypt(encryptedToken, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	cfg.Token = decryptedToken
	return &cfg, nil
}

// GetLatestConfig retrieves the most recently updated user credentials and decrypts the token
func (s *DBService) GetLatestConfig() (*GitHubConfig, error) {
	query := `SELECT username, encrypted_token, updated_at FROM github_configs ORDER BY updated_at DESC LIMIT 1`
	row := s.db.QueryRow(query)

	var cfg GitHubConfig
	var encryptedToken string
	err := row.Scan(&cfg.Username, &encryptedToken, &cfg.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	decryptedToken, err := Decrypt(encryptedToken, s.encryptionKey)
	if err != nil {
		return nil, err
	}

	cfg.Token = decryptedToken
	return &cfg, nil
}

// Encrypt encrypts standard text using AES-256-GCM and returns a hex-encoded string of nonce + ciphertext
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a hex-encoded string of nonce + ciphertext using AES-256-GCM
func Decrypt(ciphertextHex string, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
