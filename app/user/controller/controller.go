package controller

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	domain "github.com/LuchaComics/cps-backend/app/user/datastore"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// UserController Interface for user business logic controller.
type UserController interface {
	Create(ctx context.Context, m *user_s.User) (*user_s.User, error)
	CreateInitialRootAdmin(ctx context.Context) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	GetUserBySessionUUID(ctx context.Context, sessionUUID string) (*domain.User, error)
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *user_s.UserListFilter) (*user_s.UserListResult, error)
	UpdateByID(ctx context.Context, nu *user_s.User) (*user_s.User, error)
	CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*user_s.User, error)
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
