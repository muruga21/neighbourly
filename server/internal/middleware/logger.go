package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

// responseWriter is a custom wrapper to capture the HTTP status code and response body
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	body        *bytes.Buffer
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK, // Default to 200 OK
		body:           &bytes.Buffer{},
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		// Set status if Write is called without WriteHeader
		rw.WriteHeader(http.StatusOK)
	}
	// Copy the response to our buffers
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// Logger is a middleware function that logs the incoming request and the outgoing response
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture the exact request payload incoming
		var requestBodyBytes []byte
		if r.Body != nil {
			requestBodyBytes, _ = io.ReadAll(r.Body)
			// Crucial: reset the body so downstream handlers can still read it
			r.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		}

		// Wrap the ResponseWriter
		wrappedWriter := wrapResponseWriter(w)

		// Pass execution to the next handler/router
		next.ServeHTTP(wrappedWriter, r)

		// Log the complete detail of input/output
		duration := time.Since(start)
		responseBody := wrappedWriter.body.String()
		if len(requestBodyBytes) == 0 {
			requestBodyBytes = []byte("none")
		}

		log.Printf("\n==== REQUEST START ====\n"+
			"METHOD: %s\n"+
			"URL: %s\n"+
			"INPUT: %s\n"+
			"STATUS: %d\n"+
			"OUTPUT: %s\n"+
			"DURATION: %s\n"+
			"==== REQUEST END ====\n",
			r.Method, r.RequestURI,
			string(requestBodyBytes),
			wrappedWriter.status,
			responseBody,
			duration,
		)
	})
}
