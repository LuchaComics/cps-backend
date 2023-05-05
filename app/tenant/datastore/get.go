package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

func (impl TenantStorerImpl) GetByTenantID(ctx context.Context, tenantID string) (*Tenant, error) {
	filter := bson.D{{"tenant_id", tenantID}}

	var result Tenant
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by tenant id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl TenantStorerImpl) GetByName(ctx context.Context, name string) (*Tenant, error) {
	filter := bson.D{{"name", name}}

	var result Tenant
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by name error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}
