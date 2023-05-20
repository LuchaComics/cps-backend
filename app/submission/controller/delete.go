package controller

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (impl *SubmissionControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// STEP 1: Lookup the record or error.
	submission, err := impl.GetByID(ctx, id)
	if err != nil {
		log.Fatal("GetByID:", err)
	}
	if submission == nil {
		log.Fatal("GetByID: null")
	}

	// STEP 2: Delete from remote storage.
	if err := impl.S3.DeleteByKeys(ctx, []string{submission.FileUploadS3ObjectKey}); err != nil {
		impl.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
		// Continue even if we get an s3 error...
	}

	// STEP 3: Delete from database.
	if err := impl.SubmissionStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	return nil
}
