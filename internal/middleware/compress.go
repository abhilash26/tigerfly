package middleware

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

var compressibleContentTypes = map[string]bool{
	"text/html":              true,
	"text/css":               true,
	"text/plain":             true,
	"text/javascript":        true,
	"application/javascript": true,
	"application/json":       true,
	"application/atom+xml":   true,
	"application/rss+xml":    true,
	"image/svg+xml":          true,
}

func Compress(level int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if level < gzip.BestSpeed || level > gzip.BestCompression {
			level = gzip.DefaultCompression
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptEncoding := r.Header.Get("Accept-Encoding")
			if !strings.Contains(acceptEncoding, "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			contentType := w.Header().Get("Content-Type")
			if !isCompressibleContentType(contentType) {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Content-Encoding", "gzip")
			// Length is unknown after compression
			w.Header().Del("Content-Length")

			gzipWriter, err := gzip.NewWriterLevel(w, level)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			defer func() {
				if err := gzipWriter.Close(); err != nil {
					log.Printf("Error closing gzip writer: %v", err)
				}
			}()

			gzipResponse := &gzipResponseWriter{
				ResponseWriter: w,
				gzipWriter:     gzipWriter,
			}
			next.ServeHTTP(gzipResponse, r)
		})
	}
}

func isCompressibleContentType(contentType string) bool {
	for ct := range compressibleContentTypes {
		if strings.HasPrefix(contentType, ct) {
			return true
		}
	}
	return false
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func (g *gzipResponseWriter) Write(p []byte) (int, error) {
	return g.gzipWriter.Write(p)
}

// WriteHeader ensures the status code is written before any compressed data.
func (g *gzipResponseWriter) WriteHeader(statusCode int) {
	g.ResponseWriter.WriteHeader(statusCode)
}
