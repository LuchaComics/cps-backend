package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"golang.org/x/exp/slog"
)

func (c *SubmissionControllerImpl) UpdateByID(ctx context.Context, m *domain.Submission) error {
	err := c.SubmissionStorer.UpdateByID(ctx, m)
	if err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return err
}
