package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (impl *OrganizationControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	organization, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if organization == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}
	return nil
}
