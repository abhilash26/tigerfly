package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/abhilash26/tigerfly/internal/env"
)

type ServerOptions struct {
	Port              int
	IdleTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	MaxHeaderBytes    int
	ReadHeaderTimeout time.Duration
}

func New() *ServerOptions {
	defaultTimeout := 15 * time.Second

	return &ServerOptions{
		Port:              env.GetInt("SERVER_PORT", 3210),
		IdleTimeout:       env.GetDuration("SERVER_IDLE_TIMEOUT", 4*defaultTimeout),
		WriteTimeout:      env.GetDuration("SERVER_WRITE_TIMEOUT", defaultTimeout),
		ReadTimeout:       env.GetDuration("SERVER_READ_TIMEOUT", defaultTimeout),
		ReadHeaderTimeout: env.GetDuration("SERVER_READ_HEADER_TIMEOUT", defaultTimeout),
		MaxHeaderBytes:    env.GetInt("SERVER_MAX_HEADER_BYTES", 1048576),
	}
}

func (s *ServerOptions) createServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", s.Port),
		Handler:           handler,
		IdleTimeout:       s.IdleTimeout,
		WriteTimeout:      s.WriteTimeout,
		ReadTimeout:       s.ReadTimeout,
		ReadHeaderTimeout: s.ReadHeaderTimeout,
		MaxHeaderBytes:    s.MaxHeaderBytes,
	}
}

func (s *ServerOptions) Start(handler http.Handler) {
	server := s.createServer(handler)
  log.Printf("🚀 Server started on port:%d", s.Port)

	// Channel to capture OS signals for shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Graceful shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	wg.Wait()
	log.Println("Server gracefully stopped.")
}

func (s *ServerOptions) StartTLS(handler http.Handler, certFile, keyFile string) {
	server := s.createServer(handler)
	log.Printf("Server is starting with TLS on port %d", s.Port)

	// Channel to capture OS signals for shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting TLS server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down TLS server...")

	// Graceful shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during TLS server shutdown: %v", err)
	}

	wg.Wait()
	log.Println("TLS server gracefully stopped.")
}

// Fluent API setters for customization
func (s *ServerOptions) SetPort(port int) *ServerOptions {
	s.Port = port
	return s
}

func (s *ServerOptions) SetIdleTimeout(idleTimeout time.Duration) *ServerOptions {
	s.IdleTimeout = idleTimeout
	return s
}

func (s *ServerOptions) SetWriteTimeout(writeTimeout time.Duration) *ServerOptions {
	s.WriteTimeout = writeTimeout
	return s
}

func (s *ServerOptions) SetReadTimeout(readTimeout time.Duration) *ServerOptions {
	s.ReadTimeout = readTimeout
	return s
}

func (s *ServerOptions) SetMaxHeaderBytes(maxHeaderBytes int) *ServerOptions {
	s.MaxHeaderBytes = maxHeaderBytes
	return s
}

func (s *ServerOptions) SetReadHeaderTimeout(readHeaderTimeout time.Duration) *ServerOptions {
	s.ReadHeaderTimeout = readHeaderTimeout
	return s
}
