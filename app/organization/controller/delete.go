package controller

import (
	"context"

	user_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (impl *OrganizationControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		impl.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	// Update the database.
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
