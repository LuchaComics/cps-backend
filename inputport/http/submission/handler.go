package submission

import (
	submission_c "github.com/LuchaComics/cps-backend/app/submission/controller"
)

// Handler Creates http request handler
type Handler struct {
	Controller submission_c.SubmissionController
}

// NewHandler Constructor
func NewHandler(c submission_c.SubmissionController) *Handler {
	return &Handler{
		Controller: c,
	}
}
