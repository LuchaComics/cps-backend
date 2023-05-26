package customer

import (
	customer_c "github.com/LuchaComics/cps-backend/app/customer/controller"
)

// Handler Creates http request handler
type Handler struct {
	Controller customer_c.CustomerController
}

// NewHandler Constructor
func NewHandler(c customer_c.CustomerController) *Handler {
	return &Handler{
		Controller: c,
	}
}
