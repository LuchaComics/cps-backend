package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
)

func (c *CustomerControllerImpl) ListByFilter(ctx context.Context, f *user_s.UserListFilter) (*user_s.UserListResult, error) {
	// // Extract from our session the following data.
	organizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole == user_s.UserRoleRetailer {
		f.OrganizationID = organizationID
		f.Role = user_s.UserRoleCustomer
	}

	m, err := c.UserStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
