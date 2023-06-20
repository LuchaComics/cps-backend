package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	org_d "github.com/LuchaComics/cps-backend/app/organization/datastore"
	user_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
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

	// Security: Prevent deletion of root user(s).
	if organization.Type == org_d.RootType {
		impl.Logger.Warn("root organization cannot be deleted error")
		return httperror.NewForForbiddenWithSingleField("role", "root organization cannot be deleted")
	}

	// Update the database.
	organization, err := impl.GetByID(ctx, id)
	organization.Status = org_d.OrganizationArchivedStatus
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
