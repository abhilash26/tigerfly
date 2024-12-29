package middleware

import (
	"bytes"
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
			writer := &gzipResponseWriter{ResponseWriter: w, gzipWriter: gzipWriter, buffer: &bytes.Buffer{}}

			next.ServeHTTP(writer, r)

			if _, err := gzipWriter.Write(writer.buffer.Bytes()); err != nil {
				log.Printf("Error writing to gzip writer: %v", err)
			}
		})
	}
}

func isCompressibleContentType(contentType string) bool {
	if compressibleContentTypes[contentType] {
		return true
	}
	for prefix := range compressibleContentTypes {
		if strings.HasPrefix(contentType, prefix) {
			return true
		}
	}
	return false
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
	buffer     *bytes.Buffer
}

func (g *gzipResponseWriter) Write(p []byte) (int, error) {
	return g.buffer.Write(p)
}
