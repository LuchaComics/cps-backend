package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
)

func (impl UserStorerImpl) UpdateByUserID(ctx context.Context, m *User) error {
	filter := bson.D{{"user_id", m.UserID}}

	update := bson.M{ // DEVELOPERS NOTE: https://stackoverflow.com/a/60946010
		"$set": m,
	}

	// execute the UpdateOne() function to update the first matching document
	result, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update by user id error", slog.Any("error", err))
	}

	// display the number of documents updated
	impl.Logger.Debug("number of documents updated", slog.Int64("modified_count", result.ModifiedCount))

	return nil
}