package services

import (
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	key := []byte("a-very-secure-32-byte-key-123456") // 32 bytes
	plaintext := "ghp_secure_github_personal_access_token_123"

	// 1. Encrypt
	ciphertextHex, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	if len(ciphertextHex) == 0 {
		t.Fatal("ciphertext is empty")
	}

	// 2. Decrypt
	decryptedText, err := Decrypt(ciphertextHex, key)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if decryptedText != plaintext {
		t.Errorf("decrypted text mismatch: got %q, want %q", decryptedText, plaintext)
	}
}

func TestEncryptionWithInvalidKey(t *testing.T) {
	invalidKey := []byte("short-key")
	plaintext := "secret"

	_, err := Encrypt(plaintext, invalidKey)
	if err == nil {
		t.Error("expected error encrypting with short key, got nil")
	}

	_, err = Decrypt("010203040506", invalidKey)
	if err == nil {
		t.Error("expected error decrypting with short key, got nil")
	}
}
