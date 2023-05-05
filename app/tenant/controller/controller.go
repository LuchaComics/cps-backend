package controller

import (
	"golang.org/x/exp/slog"

	tenant_s "github.com/LuchaComics/cps-backend/app/tenant/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// TenantController Interface for tenant business logic controller.
type TenantController interface {
	//TODO: Add more...
}

type TenantControllerImpl struct {
	Config       *config.Conf
	Logger       *slog.Logger
	UUID         uuid.Provider
	Password     password.Provider
	TenantStorer tenant_s.TenantStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	usr_storer tenant_s.TenantStorer,
) TenantController {
	s := &TenantControllerImpl{
		Config:       appCfg,
		Logger:       loggerp,
		UUID:         uuidp,
		Password:     passwordp,
		TenantStorer: usr_storer,
	}
	s.Logger.Debug("tenant controller initialization started...")
	s.Logger.Debug("tenant controller initialized")
	return s
}
