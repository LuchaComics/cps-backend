package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	mg "github.com/LuchaComics/cps-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	s3_storage "github.com/LuchaComics/cps-backend/adapter/storage/s3"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// CustomerController Interface for customer business logic controller.
type CustomerController interface {
	Create(ctx context.Context, m *CustomerCreateRequestIDO) (*user_s.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	UpdateByID(ctx context.Context, m *user_s.User) (*user_s.User, error)
	ListByFilter(ctx context.Context, f *user_s.UserListFilter) (*user_s.UserListResult, error)
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*user_s.User, error)
}

type CustomerControllerImpl struct {
	Config      *config.Conf
	Logger      *slog.Logger
	UUID        uuid.Provider
	S3          s3_storage.S3Storager
	Password    password.Provider
	CBFFBuilder pdfbuilder.CBFFBuilder
	Emailer     mg.Emailer
	UserStorer  user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	cbffb pdfbuilder.CBFFBuilder,
	emailer mg.Emailer,
	sub_storer user_s.UserStorer,
) CustomerController {
	s := &CustomerControllerImpl{
		Config:      appCfg,
		Logger:      loggerp,
		UUID:        uuidp,
		S3:          s3,
		Password:    passwordp,
		CBFFBuilder: cbffb,
		Emailer:     emailer,
		UserStorer:  sub_storer,
	}
	s.Logger.Debug("customer controller initialization started...")
	s.Logger.Debug("customer controller initialized")
	return s
}
