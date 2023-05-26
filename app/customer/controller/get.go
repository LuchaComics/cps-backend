package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
)

func (c *CustomerControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.UserStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
