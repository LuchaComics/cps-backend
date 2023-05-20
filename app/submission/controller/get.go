package controller

import (
	"context"
	"time"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (c *SubmissionControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Submission, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.SubmissionStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}

	// The following will generate a pre-signed URL so user can download the file.
	downloadableURL, err := c.S3.GetDownloadablePresignedURL(ctx, m.FileUploadS3ObjectKey, time.Minute*15)
	if err != nil {
		c.Logger.Warn("s3 presign error", slog.Any("error", err))
		// Do not return an error, simply continue this function as there might
		// be a case were the file was removed on the s3 bucket by ourselves
		// or some other reason.
	}
	m.FileUploadDownloadableFileURL = downloadableURL

	return m, err
}
