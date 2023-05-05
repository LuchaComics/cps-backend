package controller

import (
	"context"
	"log"

	"golang.org/x/exp/slog"

	domain "github.com/LuchaComics/cps-backend/app/user/datastore"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// UserController Interface for user business logic controller.
type UserController interface {
	CreateInitialRootAdmin(ctx context.Context) error
	GetUserBySessionUUID(ctx context.Context, sessionUUID string) (*domain.User, error)
	//TODO: Add more...
}

type UserControllerImpl struct {
	Config     *config.Conf
	Logger     *slog.Logger
	UUID       uuid.Provider
	Password   password.Provider
	UserStorer user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	usr_storer user_s.UserStorer,
) UserController {
	s := &UserControllerImpl{
		Config:     appCfg,
		Logger:     loggerp,
		UUID:       uuidp,
		Password:   passwordp,
		UserStorer: usr_storer,
	}
	s.Logger.Debug("user controller initialization started...")

	// Execute the code which will check to see if we have an initial account
	// if not then we'll need to create it.
	if err := s.CreateInitialRootAdmin(context.Background()); err != nil {
		log.Fatal(err) // We terminate app here b/c dependency injection not allowed to fail, so fail here at startup of dynamodb.
	}

	s.Logger.Debug("user controller initialized")
	return s
}
