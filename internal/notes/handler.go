package notes

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/kalpesh122/agentic-backend-go/internal/apperr"
	"github.com/kalpesh122/agentic-backend-go/internal/auth"
	"github.com/kalpesh122/agentic-backend-go/internal/httpx/respond"
)

const maxBodyBytes = 64 << 10

// Handler exposes the notes service over HTTP. It only translates HTTP ↔ service types.
type Handler struct{ svc *Service }

// NewHandler wires the service.
func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// Routes mounts the resource on a chi router. Auth is applied by the caller.
func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Patch("/{id}", h.update)
	r.Delete("/{id}", h.remove)
}

func parseID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, apperr.Validation(map[string]string{"id": "must be a UUID"})
	}
	return id, nil
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			respond.Error(w, r, apperr.Validation(map[string]string{"limit": "must be an integer"}))
			return
		}
		limit = n
	}
	page, err := h.svc.List(r.Context(), auth.Principal(r.Context()), limit, r.URL.Query().Get("cursor"))
	if err != nil {
		respond.Error(w, r, err)
		return
	}
	respond.JSON(w, http.StatusOK, page)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in CreateInput
	if err := respond.Decode(r, &in, maxBodyBytes); err != nil {
		respond.Error(w, r, err)
		return
	}
	n, err := h.svc.Create(r.Context(), auth.Principal(r.Context()), in)
	if err != nil {
		respond.Error(w, r, err)
		return
	}
	respond.JSON(w, http.StatusCreated, n)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respond.Error(w, r, err)
		return
	}
	n, err := h.svc.Get(r.Context(), auth.Principal(r.Context()), id)
	if err != nil {
		respond.Error(w, r, err)
		return
	}
	respond.JSON(w, http.StatusOK, n)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respond.Error(w, r, err)
		return
	}
	var in UpdateInput
	if err := respond.Decode(r, &in, maxBodyBytes); err != nil {
		respond.Error(w, r, err)
		return
	}
	n, err := h.svc.Update(r.Context(), auth.Principal(r.Context()), id, in)
	if err != nil {
		respond.Error(w, r, err)
		return
	}
	respond.JSON(w, http.StatusOK, n)
}

func (h *Handler) remove(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respond.Error(w, r, err)
		return
	}
	if err := h.svc.Delete(r.Context(), auth.Principal(r.Context()), id); err != nil {
		respond.Error(w, r, err)
		return
	}
	respond.JSON(w, http.StatusNoContent, nil)
}
