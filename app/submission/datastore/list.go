package datastore

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl SubmissionStorerImpl) ListByFilter(ctx context.Context, f *SubmissionListFilter) (*SubmissionListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// Get a reference to the collection
	collection := impl.Collection

	// Pagination parameters
	pageSize := 10
	startAfter := "" // The ID to start after, initially empty for the first page

	// Sorting parameters
	sortField := "_id"
	sortOrder := 1 // 1=ascending | -1=descending

	// Pagination filter
	filter := bson.M{}
	options := options.Find().
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{sortField, sortOrder}})

	// Add filter conditions to the filter
	if f.UserID != primitive.NilObjectID {
		filter["user_id"] = f.UserID
	}
	if f.UserRole != 0 {
		filter["user_role"] = f.UserRole
	}
	if f.OrganizationID != primitive.NilObjectID {
		filter["organization_id"] = f.OrganizationID
	}

	if startAfter != "" {
		// Find the document with the given startAfter ID
		cursor, err := collection.FindOne(ctx, bson.M{"_id": startAfter}).DecodeBytes()
		if err != nil {
			log.Fatal(err)
		}
		options.SetSkip(1)
		filter["_id"] = bson.M{"$gt": cursor.Lookup("_id").ObjectID()}
	}

	if f.ExcludeArchived {
		filter["state"] = bson.M{"$ne": SubmissionArchivedState} // Do not list archived items! This code
	}

	options.SetSort(bson.D{{sortField, 1}}) // Sort in ascending order based on the specified field

	// Retrieve the list of items from the collection
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results = []*Submission{}
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	return &SubmissionListResult{
		Results: results,
	}, nil
}
