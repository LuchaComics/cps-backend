package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (impl SubmissionStorerImpl) CountByFilter(ctx context.Context, f *SubmissionListFilter) (int64, error) {

	filter := bson.M{}

	if f.OrganizationID != primitive.NilObjectID {
		filter["organization_id"] = f.OrganizationID
	}

	if f.UserID != primitive.NilObjectID {
		filter["user_id"] = f.UserID
	}

	if f.UserRole != 0 {
		filter["user_role"] = f.UserRole
	}

	if f.ExcludeArchived {
		filter["state"] = bson.M{"$ne": SubmissionArchivedState} // Do not list archived items! This code
	}

	opts := options.Count().SetHint("_id_")
	count, err := impl.Collection.CountDocuments(ctx, filter, opts)
	if err != nil {
		return 0, err
	}

	return count, nil
}
