package handler

import (
	"github.com/ronaudinho/bigq/internal/service"

	"github.com/labstack/echo/v4"
)

// Handler handles communication over HTTP
type Handler struct {
	service *service.Service
}

// New creates new instance of Handler
func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Send picks up queued results from exhange and queue
// NOTE it is probably possible to combine all into one
// using switch statement and interface{} inside
func (h *Handler) Send(c echo.Context) {}
