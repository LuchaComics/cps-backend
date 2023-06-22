package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (impl *ComicSubmissionControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// STEP 1: Lookup the record or error.
	submission, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if submission == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}

	// STEP 2: Delete from remote storage.
	if err := impl.S3.DeleteByKeys(ctx, []string{submission.FileUploadS3ObjectKey}); err != nil {
		impl.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
		// Do not return an error, simply continue this function as there might
		// be a case were the file was removed on the s3 bucket by ourselves
		// or some other reason.
	}

	// STEP 3: Delete from database.
	if err := impl.ComicSubmissionStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	return nil
}
