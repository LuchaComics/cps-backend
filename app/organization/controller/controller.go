package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	mg "github.com/LuchaComics/cps-backend/adapter/emailer/mailgun"
	s3_storage "github.com/LuchaComics/cps-backend/adapter/storage/s3"
	comicsub_s "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	domain "github.com/LuchaComics/cps-backend/app/organization/datastore"
	org_d "github.com/LuchaComics/cps-backend/app/organization/datastore"
	organization_s "github.com/LuchaComics/cps-backend/app/organization/datastore"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// OrganizationController Interface for organization business logic controller.
type OrganizationController interface {
	Create(ctx context.Context, m *domain.Organization) (*domain.Organization, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Organization, error)
	UpdateByID(ctx context.Context, m *domain.Organization) (*domain.Organization, error)
	ListByFilter(ctx context.Context, f *domain.OrganizationListFilter) (*domain.OrganizationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.OrganizationListFilter) ([]*domain.OrganizationAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*org_d.Organization, error)
}

type OrganizationControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	S3                    s3_storage.S3Storager
	Emailer               mg.Emailer
	OrganizationStorer    organization_s.OrganizationStorer
	UserStorer            user_s.UserStorer
	ComicSubmissionStorer comicsub_s.ComicSubmissionStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	emailer mg.Emailer,
	org_storer organization_s.OrganizationStorer,
	usr_storer user_s.UserStorer,
	csub_storer comicsub_s.ComicSubmissionStorer,
) OrganizationController {
	s := &OrganizationControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		S3:                    s3,
		Emailer:               emailer,
		OrganizationStorer:    org_storer,
		UserStorer:            usr_storer,
		ComicSubmissionStorer: csub_storer,
	}
	s.Logger.Debug("organization controller initialization started...")
	s.Logger.Debug("organization controller initialized")
	return s
}
