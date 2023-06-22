package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	organization_s "github.com/LuchaComics/cps-backend/app/organization/datastore"
	domain "github.com/LuchaComics/cps-backend/app/user/datastore"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// UserController Interface for user business logic controller.
type UserController interface {
	Create(ctx context.Context, requestData *UserCreateRequestIDO) (*user_s.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	GetUserBySessionUUID(ctx context.Context, sessionUUID string) (*domain.User, error)
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *user_s.UserListFilter) (*user_s.UserListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *user_s.UserListFilter) ([]*user_s.UserAsSelectOption, error)
	UpdateByID(ctx context.Context, request *UserUpdateRequestIDO) (*user_s.User, error)
	CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*user_s.User, error)
	//TODO: Add more...
}

type UserControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	Password           password.Provider
	OrganizationStorer organization_s.OrganizationStorer
	UserStorer         user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	org_storer organization_s.OrganizationStorer,
	usr_storer user_s.UserStorer,
) UserController {
	s := &UserControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		Password:           passwordp,
		OrganizationStorer: org_storer,
		UserStorer:         usr_storer,
	}
	s.Logger.Debug("user controller initialization started...")

	s.Logger.Debug("user controller initialized")
	return s
}
