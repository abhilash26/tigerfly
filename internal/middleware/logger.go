package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Buffered channel for async logging
var logChannel = make(chan string, 100)

func init() {
	go func() {
		for logEntry := range logChannel {
			log.Println(logEntry)
		}
	}()
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &WrapResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		start := time.Now()

		// Serve the request
		next.ServeHTTP(ww, r)

		// Log the request asynchronously
		logRequest(r, ww.statusCode, ww.bytes, time.Since(start))
	})
}

type WrapResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (w *WrapResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *WrapResponseWriter) Write(p []byte) (int, error) {
	bytesWritten, err := w.ResponseWriter.Write(p)
	w.bytes += bytesWritten
	return bytesWritten, err
}

func logRequest(r *http.Request, status int, size int, duration time.Duration) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	var statusColor string
	switch {
	case status >= 100 && status <= 199:
		statusColor = "\033[34m"
	case status >= 200 && status <= 299:
		statusColor = "\033[32m"
	case status >= 300 && status <= 399:
		statusColor = "\033[33m"
	case status >= 400 && status <= 599:
		statusColor = "\033[31m"
	default:
		statusColor = "\033[0m"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s%03d\033[0m", statusColor, status))                          // Status code
	builder.WriteString(fmt.Sprintf(" \033[35m%s\033[36m %s\033[0m", r.Method, fullURL))            // Method and URL
	builder.WriteString(fmt.Sprintf(" \033[33m%dB\033[0m", size))                                   // Response size
	builder.WriteString(fmt.Sprintf(" in \033[34m%.1fµs\033[0m", float64(duration.Microseconds()))) // Duration

	// Send log entry to the async channel
	logChannel <- builder.String()
}
