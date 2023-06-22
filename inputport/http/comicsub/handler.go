package comicsub

import (
	comicsub_c "github.com/LuchaComics/cps-backend/app/comicsub/controller"
)

// Handler Creates http request handler
type Handler struct {
	Controller comicsub_c.ComicSubmissionController
}

// NewHandler Constructor
func NewHandler(c comicsub_c.ComicSubmissionController) *Handler {
	return &Handler{
		Controller: c,
	}
}
