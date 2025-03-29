package network

import (
	"NotRedis/internal/compute/parser"
	"NotRedis/internal/database/engine"
	"NotRedis/internal/database/storage"
	"bufio"
	"go.uber.org/zap"
	"net"
	"testing"
	"time"
)

func setupTestServer() (*Server, *zap.Logger) {
	logger, _ := zap.NewDevelopment()
	cfg := Config{
		Address:        "127.0.0.1:0",
		MaxConnections: 2,
		MaxMessageSize: 1024,
		IdleTimeout:    1 * time.Second,
	}
	myMap := engine.NewMyMap(100, logger)
	store, _ := storage.NewStorage(myMap, logger)
	p := parser.NewParser(logger)
	srv := NewServer(cfg, store, p, logger)
	return srv, logger
}

func TestServer_Start(t *testing.T) {
	srv, _ := setupTestServer()

	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", srv.config.Address)
	if err != nil {
		t.Errorf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	if conn == nil {
		t.Error("Client connection should not be nil")
	}
}

func TestServer_MaxConnections(t *testing.T) {
	srv, _ := setupTestServer()

	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	conn1, err := net.Dial("tcp", srv.config.Address)
	if err != nil {
		t.Errorf("Failed to connect first client: %v", err)
	}
	defer conn1.Close()

	conn2, err := net.Dial("tcp", srv.config.Address)
	if err != nil {
		t.Errorf("Failed to connect second client: %v", err)
	}
	defer conn2.Close()

	conn3, err := net.Dial("tcp", srv.config.Address)
	if err != nil {
		t.Errorf("Failed to connect third client: %v", err)
	}
	defer conn3.Close()

	reader := bufio.NewReader(conn3)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("Failed to read response from third client: %v", err)
	}
	if response != "Too many connections\n" {
		t.Errorf("Expected 'Too many connections', got %q", response)
	}
}

func TestServer_HandleConnection(t *testing.T) {
	srv, _ := setupTestServer()

	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", srv.config.Address)
	if err != nil {
		t.Errorf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	_, err = writer.WriteString("SET key value\n")
	if err != nil {
		t.Errorf("Failed to write SET command: %v", err)
	}
	writer.Flush()
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("Failed to read SET response: %v", err)
	}
	if response != "OK\n" {
		t.Errorf("Expected 'OK' for SET, got %q", response)
	}

	_, err = writer.WriteString("GET key\n")
	if err != nil {
		t.Errorf("Failed to write GET command: %v", err)
	}
	writer.Flush()
	response, err = reader.ReadString('\n')
	if err != nil {
		t.Errorf("Failed to read GET response: %v", err)
	}
	if response != "value\n" {
		t.Errorf("Expected 'value' for GET, got %q", response)
	}

	_, err = writer.WriteString("DEL key\n")
	if err != nil {
		t.Errorf("Failed to write DEL command: %v", err)
	}
	writer.Flush()
	response, err = reader.ReadString('\n')
	if err != nil {
		t.Errorf("Failed to read DEL response: %v", err)
	}
	if response != "OK\n" {
		t.Errorf("Expected 'OK' for DEL, got %q", response)
	}

	_, err = writer.WriteString("GET key\n")
	if err != nil {
		t.Errorf("Failed to write GET command after DEL: %v", err)
	}
	writer.Flush()
	response, err = reader.ReadString('\n')
	if err != nil {
		t.Errorf("Failed to read GET response after DEL: %v", err)
	}
	if !containsError(response) {
		t.Errorf("Expected error for GET after DEL, got %q", response)
	}
}

func containsError(s string) bool {
	return len(s) > 5 && s[:5] == "Error"
}

func TestServer_IdleTimeout(t *testing.T) {
	srv, _ := setupTestServer()

	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", srv.config.Address)
	if err != nil {
		t.Errorf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	time.Sleep(2 * time.Second)

	_, err = writer.WriteString("GET key\n")
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}

	_, err = reader.ReadString('\n')
	if err == nil {
		t.Error("Expected error due to idle timeout, but read succeeded")
	}
}
