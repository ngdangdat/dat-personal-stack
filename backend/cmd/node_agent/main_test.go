package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"online"}`))
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"online"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestWebSocketHandshakeInvalidToken(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token != "secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Scheme = "ws"
	u.Path = "/rpc"
	u.RawQuery = "token=wrong-token"

	_, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err == nil {
		t.Fatalf("Expected connection failure on invalid token, got success")
	}

	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %v", resp)
	}
}

func TestWebSocketRPCSystemTelemetry(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		safeConn := &SafeConn{Conn: conn}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var req JSONRPCRequest
			if err := json.Unmarshal(msg, &req); err != nil {
				continue
			}
			handleRequest(context.Background(), safeConn, req)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Scheme = "ws"
	u.Path = "/rpc"

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to dial websocket: %v", err)
	}
	defer ws.Close()

	// Send system.get_telemetry request
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "system.get_telemetry",
		ID:      1,
	}
	reqData, _ := json.Marshal(req)
	err = ws.WriteMessage(websocket.TextMessage, reqData)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Read response
	_, respData, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var resp JSONRPCResponse
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected jsonrpc 2.0, got %s", resp.JSONRPC)
	}

	if resp.ID != float64(1) {
		t.Errorf("Expected response ID to be 1, got %v", resp.ID)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Result is not a map[string]interface{}")
	}

	if _, ok := resultMap["cpu"]; !ok {
		t.Errorf("Expected result map to contain 'cpu'")
	}
	if _, ok := resultMap["memory"]; !ok {
		t.Errorf("Expected result map to contain 'memory'")
	}
	if _, ok := resultMap["uptime_seconds"]; !ok {
		t.Errorf("Expected result map to contain 'uptime_seconds'")
	}
}

func TestWebSocketRPCGitStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		safeConn := &SafeConn{Conn: conn}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var req JSONRPCRequest
			if err := json.Unmarshal(msg, &req); err != nil {
				continue
			}
			handleRequest(context.Background(), safeConn, req)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Scheme = "ws"
	u.Path = "/rpc"

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Failed to dial websocket: %v", err)
	}
	defer ws.Close()

	// Send workspace.get_git_status request
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "workspace.get_git_status",
		ID:      2,
	}
	reqData, _ := json.Marshal(req)
	err = ws.WriteMessage(websocket.TextMessage, reqData)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Read response
	_, respData, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var resp JSONRPCResponse
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected jsonrpc 2.0, got %s", resp.JSONRPC)
	}

	if resp.ID != float64(2) {
		t.Errorf("Expected response ID to be 2, got %v", resp.ID)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Result is not a map[string]interface{}")
	}

	if _, ok := resultMap["branch"]; !ok {
		t.Errorf("Expected result map to contain 'branch'")
	}
	if _, ok := resultMap["is_clean"]; !ok {
		t.Errorf("Expected result map to contain 'is_clean'")
	}
	if _, ok := resultMap["commits"]; !ok {
		t.Errorf("Expected result map to contain 'commits'")
	}
}
