package controller

import (
	"context"

	"golang.org/x/exp/slog"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	submission_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// SubmissionController Interface for submission business logic controller.
type SubmissionController interface {
	Create(ctx context.Context, m *domain.Submission) error
	GetBySubmissionID(ctx context.Context, submissionID string) (*domain.Submission, error)
	UpdateBySubmissionID(ctx context.Context, m *domain.Submission) error
}

type SubmissionControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	Password         password.Provider
	SubmissionStorer submission_s.SubmissionStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	sub_storer submission_s.SubmissionStorer,
) SubmissionController {
	s := &SubmissionControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		Password:         passwordp,
		SubmissionStorer: sub_storer,
	}
	s.Logger.Debug("submission controller initialization started...")
	s.Logger.Debug("submission controller initialized")
	return s
}
