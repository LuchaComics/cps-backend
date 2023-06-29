package attachment

import (
	attachment_c "github.com/LuchaComics/cps-backend/app/attachment/controller"
)

// Handler Creates http request handler
type Handler struct {
	Controller attachment_c.AttachmentController
}

// NewHandler Constructor
func NewHandler(c attachment_c.AttachmentController) *Handler {
	return &Handler{
		Controller: c,
	}
}
