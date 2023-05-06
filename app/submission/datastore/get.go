package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

func (impl SubmissionStorerImpl) GetBySubmissionID(ctx context.Context, submissionID string) (*Submission, error) {
	filter := bson.D{{"submission_id", submissionID}}

	var result Submission
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by submission id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}
