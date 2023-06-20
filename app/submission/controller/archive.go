package controller

import (
	"context"
	"time"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (c *SubmissionControllerImpl) ArchiveByID(ctx context.Context, id primitive.ObjectID) (*domain.Submission, error) {
	// Fetch the original submission.
	os, err := c.SubmissionStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, nil
	}

	// Modify our original submission.
	os.ModifiedAt = time.Now()
	os.Status = domain.SubmissionStatusArchived

	// Save to the database the modified submission.
	if err := c.SubmissionStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}
