package tenant

import (
	tenant_c "github.com/LuchaComics/cps-backend/app/tenant/controller"
)

// Handler Creates http request handler
type Handler struct {
	Controller tenant_c.TenantController
}

// NewHandler Constructor
func NewHandler(c tenant_c.TenantController) *Handler {
	return &Handler{
		Controller: c,
	}
}
