package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
)

type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	RemoteIP  string    `json:"remote_ip"`
	Payload   string    `json:"payload,omitempty"`
}

type AuditMiddleware struct {
	nc *nats.Conn
}

func NewAuditMiddleware(nc *nats.Conn) *AuditMiddleware {
	return &AuditMiddleware{nc: nc}
}

func (m *AuditMiddleware) Audit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only audit write operations
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			var body []byte
			if r.Body != nil {
				body, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			event := AuditEvent{
				Timestamp: time.Now(),
				Method:    r.Method,
				Path:      r.URL.Path,
				RemoteIP:  r.RemoteAddr,
				Payload:   string(body),
			}

			data, _ := json.Marshal(event)
			if m.nc != nil {
				m.nc.Publish("nebula.audit", data)
			}
		}

		next.ServeHTTP(w, r)
	})
}
