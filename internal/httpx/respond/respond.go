// Package respond writes JSON responses and the standard error body.
package respond

import (
	"encoding/json/v2"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/kalpesh122/agentic-backend-go/internal/apperr"
)

// ErrorBody is the wire format for every error: {"error": {code, message, details?, request_id}}.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail carries the stable code and a request id for correlation.
type ErrorDetail struct {
	Code      apperr.Code `json:"code"`
	Message   string      `json:"message"`
	Details   any         `json:"details,omitempty"`
	RequestID string      `json:"request_id"`
}

// JSON writes v with the given status. A nil v with 204 writes no body.
func JSON(w http.ResponseWriter, status int, v any) {
	if v == nil {
		w.WriteHeader(status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.MarshalWrite(w, v); err != nil {
		slog.Error("write response", "err", err)
	}
}

// Error maps any error to the standard body. Unknown errors become 500 and are logged with the request id.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	rid := middleware.GetReqID(r.Context())
	if e, ok := apperr.As(err); ok {
		if e.Status >= 500 {
			slog.ErrorContext(r.Context(), "request failed", "err", err, "request_id", rid)
		}
		JSON(w, e.Status, ErrorBody{Error: ErrorDetail{Code: e.Code, Message: e.Message, Details: e.Details, RequestID: rid}})
		return
	}
	slog.ErrorContext(r.Context(), "unhandled error", "err", err, "request_id", rid)
	JSON(w, http.StatusInternalServerError, ErrorBody{Error: ErrorDetail{Code: apperr.CodeInternal, Message: "Internal server error", RequestID: rid}})
}

// Decode reads a JSON body into v, rejecting unknown fields and oversized payloads.
func Decode(r *http.Request, v any, maxBytes int64) error {
	r.Body = http.MaxBytesReader(nil, r.Body, maxBytes)
	if err := json.UnmarshalRead(r.Body, v, json.RejectUnknownMembers(true)); err != nil {
		return apperr.Wrap(apperr.Validation(map[string]string{"body": "malformed JSON: " + err.Error()}), err)
	}
	return nil
}
