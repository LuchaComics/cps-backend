package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

func (impl SubmissionStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*Submission, error) {
	filter := bson.M{"_id": id}

	var result Submission
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl SubmissionStorerImpl) GetByCPSRN(ctx context.Context, cpsrn string) (*Submission, error) {
	filter := bson.M{"cpsrn": cpsrn}

	var result Submission
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}
