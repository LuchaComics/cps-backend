package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/organization/datastore"
	user_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (c *OrganizationControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Organization, error) {
	// Extract from our session the following data.
	userOrganizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// If user is not administrator nor belongs to the organization then error.
	if userRole != user_d.UserRoleRoot && id != userOrganizationID {
		c.Logger.Error("authenticated user is not staff role nor belongs to the organization error",
			slog.Any("userRole", userRole),
			slog.Any("userOrganizationID", userOrganizationID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this organization")
	}

	// Retrieve from our database the record for the specific id.
	m, err := c.OrganizationStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
