package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"golang.org/x/exp/slog"
)

func (c *SubmissionControllerImpl) UpdateBySubmissionID(ctx context.Context, m *domain.Submission) error {
	err := c.SubmissionStorer.UpdateBySubmissionID(ctx, m)
	if err != nil {
		c.Logger.Error("database update by submission id error", slog.Any("error", err))
		return err
	}
	return err
}
