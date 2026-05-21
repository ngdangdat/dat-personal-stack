package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Config holds agent parameters
type Config struct {
	Port        string
	SecretToken string
}

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// RPCError represents a JSON-RPC 2.0 error object
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RunCommandParams parameters for workspace.run_command
type RunCommandParams struct {
	WorkspaceID string `json:"workspace_id"`
	Command     string `json:"command"`
}

// StreamNotification represents a custom notification (stdout/stderr)
type StreamNotification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type OutputParams struct {
	WorkspaceID string      `json:"workspace_id"`
	Data        string      `json:"data"`
	CommandID   interface{} `json:"command_id"`
}

var (
	startTime = time.Now()
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins for simplicity in local development
			return true
		},
	}
)

type SafeConn struct {
	*websocket.Conn
	mu sync.Mutex
}

func (sc *SafeConn) WriteJSONSafe(v interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.Conn.WriteJSON(v)
}

func main() {
	config := Config{
		Port:        getEnv("PORT", "8081"),
		SecretToken: getEnv("AGENT_SECRET_TOKEN", "agent-secret-token"),
	}

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"online"}`))
	})

	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		// Handshake token validation
		token := r.URL.Query().Get("token")
		if token != config.SecretToken {
			http.Error(w, "Unauthorized: Invalid secret token", http.StatusUnauthorized)
			log.Printf("Rejected connection: Invalid token '%s'", token)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}
		defer conn.Close()

		safeConn := &SafeConn{Conn: conn}
		log.Println("Client connected successfully via WebSocket")

		// Start periodic telemetry sender (every 3 seconds)
		stopTelemetry := make(chan struct{})
		go startTelemetryStream(safeConn, stopTelemetry)
		defer close(stopTelemetry)

		// Read loop
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read error: %v", err)
				break
			}

			var req JSONRPCRequest
			if err := json.Unmarshal(msg, &req); err != nil {
				sendError(safeConn, -32700, "Parse error", nil, nil)
				continue
			}

			handleRequest(safeConn, req)
		}
		log.Println("Client disconnected")
	})

	log.Printf("Starting node agent daemon on :%s...", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Agent failed to start: %v", err)
	}
}

func handleRequest(conn *SafeConn, req JSONRPCRequest) {
	if req.JSONRPC != "2.0" {
		sendError(conn, -32600, "Invalid Request: expected jsonrpc version '2.0'", nil, req.ID)
		return
	}

	switch req.Method {
	case "system.get_telemetry":
		telemetry := getSystemTelemetry()
		conn.WriteJSONSafe(JSONRPCResponse{
			JSONRPC: "2.0",
			Result:  telemetry,
			ID:      req.ID,
		})

	case "workspace.run_command":
		var params RunCommandParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			sendError(conn, -32602, "Invalid params", nil, req.ID)
			return
		}
		go executeCommand(conn, params, req.ID)

	default:
		sendError(conn, -32601, fmt.Sprintf("Method not found: '%s'", req.Method), nil, req.ID)
	}
}

func getSystemTelemetry() map[string]interface{} {
	// Calculate simulated CPU telemetry with slight noise
	cpuPercent := 10.0 + rand.Float64()*15.0 // baseline 10-25%
	temp := 40.0 + rand.Float64()*10.0       // baseline 40-50 C

	// Goroutine/Thread counts
	threads := runtime.NumGoroutine()

	// Memory telemetry (simulating 8GB system)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	totalBytes := uint64(8589934592)
	usedBytes := uint64(2000000000) + m.Alloc

	return map[string]interface{}{
		"cpu": map[string]interface{}{
			"usage_percent": cpuPercent,
			"temp_celsius":  temp,
		},
		"memory": map[string]interface{}{
			"total_bytes": totalBytes,
			"used_bytes":  usedBytes,
		},
		"uptime_seconds": int64(time.Since(startTime).Seconds()),
		"threads":        threads,
	}
}

func startTelemetryStream(conn *SafeConn, stop chan struct{}) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			telemetry := getSystemTelemetry()
			// Stream telemetry updates as notifications (no ID)
			err := conn.WriteJSONSafe(StreamNotification{
				JSONRPC: "2.0",
				Method:  "system.telemetry_stream",
				Params:  telemetry,
			})
			if err != nil {
				return
			}
		case <-stop:
			return
		}
	}
}

func executeCommand(conn *SafeConn, params RunCommandParams, rpcID interface{}) {
	cmdStr := params.Command
	if cmdStr == "" {
		sendError(conn, -32602, "Missing command string", nil, rpcID)
		return
	}

	log.Printf("Executing command: %s (Workspace: %s)", cmdStr, params.WorkspaceID)

	// In container/local environments, run in a shell context
	cmd := exec.Command("sh", "-c", cmdStr)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		sendError(conn, -32000, fmt.Sprintf("Failed to open stdout: %v", err), nil, rpcID)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		sendError(conn, -32000, fmt.Sprintf("Failed to open stderr: %v", err), nil, rpcID)
		return
	}

	if err := cmd.Start(); err != nil {
		sendError(conn, -32000, fmt.Sprintf("Failed to start process: %v", err), nil, rpcID)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Stream stdout chunks
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdoutPipe)
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				conn.WriteJSONSafe(StreamNotification{
					JSONRPC: "2.0",
					Method:  "workspace.stdout",
					Params: OutputParams{
						WorkspaceID: params.WorkspaceID,
						Data:        string(buf[:n]),
						CommandID:   rpcID,
					},
				})
			}
			if err != nil {
				break
			}
		}
	}()

	// Stream stderr chunks
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderrPipe)
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				conn.WriteJSONSafe(StreamNotification{
					JSONRPC: "2.0",
					Method:  "workspace.stderr",
					Params: OutputParams{
						WorkspaceID: params.WorkspaceID,
						Data:        string(buf[:n]),
						CommandID:   rpcID,
					},
				})
			}
			if err != nil {
				break
			}
		}
	}()

	wg.Wait()
	exitCode := 0
	if err := cmd.Wait(); err != nil {
		exitCode = 1
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	// Send final RPC completion response
	conn.WriteJSONSafe(JSONRPCResponse{
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"exit_code": exitCode,
			"message":   "Command execution finished",
		},
		ID: rpcID,
	})
}

func sendError(conn *SafeConn, code int, message string, data interface{}, id interface{}) {
	conn.WriteJSONSafe(JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	})
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
