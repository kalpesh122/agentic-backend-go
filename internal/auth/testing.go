package auth

import "github.com/kalpesh122/agentic-backend-go/internal/apperr"

// ErrTestUnauthorized mirrors the production 401 so handler tests can assert the same body.
func ErrTestUnauthorized() *apperr.Error { return apperr.Unauthorized("") }
