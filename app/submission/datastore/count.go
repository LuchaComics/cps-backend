package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl SubmissionStorerImpl) CountAll(ctx context.Context) (int64, error) {

	opts := options.Count().SetHint("_id_")
	count, err := impl.Collection.CountDocuments(ctx, bson.D{}, opts)
	if err != nil {
		return 0, err
	}

	return count, nil
}
