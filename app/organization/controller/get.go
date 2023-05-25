package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/organization/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (c *OrganizationControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Organization, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.OrganizationStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
