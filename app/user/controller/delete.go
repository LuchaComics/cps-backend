package controller

import (
	"context"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (impl *UserControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// Extract from our session the following data.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole != user_s.UserRoleRoot {
		return httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// STEP 1: Lookup the record or error.
	user, err := impl.UserStorer.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if user == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}

	// STEP 2: Delete from database.
	if err := impl.UserStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database delete by id error", slog.Any("error", err))
		return err
	}
	return nil
}
