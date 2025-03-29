package network

import (
	"NotRedis/internal/compute/parser"
	"NotRedis/internal/database/storage"
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

type Config struct {
	Address        string
	MaxConnections int
	MaxMessageSize int
	IdleTimeout    time.Duration
}

type Server struct {
	config      Config
	store       *storage.Storage
	parser      *parser.Parser
	logger      *zap.Logger
	mu          sync.Mutex
	connections int
}

func NewServer(cfg Config, store *storage.Storage, parser *parser.Parser, logger *zap.Logger) *Server {
	return &Server{
		config:      cfg,
		store:       store,
		parser:      parser,
		logger:      logger,
		connections: 0,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		s.logger.Error("failed to start server")
		return err
	}
	defer listener.Close()

	s.config.Address = listener.Addr().String()
	s.logger.Info("server started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("failed to accept connection")
			continue
		}

		if s.addConnection() {
			go s.handleConnection(conn)
		} else {
			s.logger.Error("сonnection limit exceeded")
			conn.Write([]byte("Too many connections\n"))
			conn.Close()
		}
	}
}

func (s *Server) addConnection() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.connections >= s.config.MaxConnections {
		return false
	}
	s.connections++
	return true
}

func (s *Server) removeConnection() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections--
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.removeConnection()
		conn.Close()
	}()

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic in connection handler")
		}
	}()

	s.logger.Info("new client connected")

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		conn.SetDeadline(time.Now().Add(s.config.IdleTimeout))

		request, err := reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				s.logger.Info("connection timed out")
			} else {
				s.logger.Error("failed to read request")
			}
			return
		}

		if len(request) > s.config.MaxMessageSize {
			writer.WriteString("Message too large\n")
			writer.Flush()
			continue
		}

		req, err := s.parser.Parse(request)
		if err != nil {
			writer.WriteString(fmt.Sprintf("Error: %v\n", err))
			writer.Flush()
			continue
		}

		switch req.Type {
		case "SET":
			err = s.store.Set(req.Key, req.Value)
			if err != nil {
				writer.WriteString(fmt.Sprintf("Error: %v\n", err))
			} else {
				writer.WriteString("OK\n")
			}
		case "GET":
			value, err := s.store.Get(req.Key)
			if err != nil {
				writer.WriteString(fmt.Sprintf("Error: %v\n", err))
			} else {
				writer.WriteString(fmt.Sprintf("%s\n", value))
			}
		case "DEL":
			err = s.store.Del(req.Key)
			if err != nil {
				writer.WriteString(fmt.Sprintf("Error: %v\n", err))
			} else {
				writer.WriteString("OK\n")
			}
		}
		writer.Flush()
	}
}
