package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"golang.org/x/exp/slog"
)

func (c *SubmissionControllerImpl) GetBySubmissionID(ctx context.Context, submissionID string) (*domain.Submission, error) {
	m, err := c.SubmissionStorer.GetBySubmissionID(ctx, submissionID)
	if err != nil {
		c.Logger.Error("database get by submission id error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
