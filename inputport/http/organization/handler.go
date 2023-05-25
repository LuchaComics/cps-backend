package organization

import (
	organization_c "github.com/LuchaComics/cps-backend/app/organization/controller"
)

// Handler Creates http request handler
type Handler struct {
	Controller organization_c.OrganizationController
}

// NewHandler Constructor
func NewHandler(c organization_c.OrganizationController) *Handler {
	return &Handler{
		Controller: c,
	}
}
