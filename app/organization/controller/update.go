package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	domain "github.com/LuchaComics/cps-backend/app/organization/datastore"
	user_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (c *OrganizationControllerImpl) UpdateByID(ctx context.Context, ns *domain.Organization) (*domain.Organization, error) {
	// Fetch the original organization.
	os, err := c.OrganizationStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		c.Logger.Error("organization does not exist error",
			slog.Any("organization_id", ns.ID))
		return nil, httperror.NewForBadRequestWithSingleField("message", "organization does not exist")
	}

	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOrganizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)
	userName := ctx.Value(constants.SessionUserName).(string)

	// If user is not administrator nor belongs to the organization then error.
	if userRole != user_d.StaffRole && os.ID != userOrganizationID {
		c.Logger.Error("authenticated user is not staff role nor belongs to the organization error",
			slog.Any("userRole", userRole),
			slog.Any("userOrganizationID", userOrganizationID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this organization")
	}

	// Modify our original organization.
	os.ModifiedAt = time.Now()
	os.ModifiedByUserID = userID
	os.ModifiedByUserName = userName
	os.Type = ns.Type
	os.State = ns.State
	os.Name = ns.Name

	// Save to the database the modified organization.
	if err := c.OrganizationStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}
