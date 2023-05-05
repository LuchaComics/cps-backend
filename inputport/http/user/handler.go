package user

import (
	user_c "github.com/LuchaComics/cps-backend/app/user/controller"
)

// Handler Creates http request handler
type Handler struct {
	Controller user_c.UserController
}

// NewHandler Constructor
func NewHandler(c user_c.UserController) *Handler {
	return &Handler{
		Controller: c,
	}
}
