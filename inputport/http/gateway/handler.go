package gateway

import (
	gateway_c "github.com/LuchaComics/cps-backend/app/gateway/controller"
)

// Handler Creates http request handler
type Handler struct {
	Controller gateway_c.GatewayController
}

// NewHandler Constructor
func NewHandler(c gateway_c.GatewayController) *Handler {
	return &Handler{
		Controller: c,
	}
}
