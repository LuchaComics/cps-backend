package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
)

func (c *SubmissionControllerImpl) GetSubmissionBySessionUUID(ctx context.Context, sessionUUID string) (*domain.Submission, error) {
	panic("TODO: IMPLEMENT")
	return nil, nil
}
