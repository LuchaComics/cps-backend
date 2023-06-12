package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	s_d "github.com/LuchaComics/cps-backend/app/organization/datastore"
	user_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (c *OrganizationControllerImpl) Create(ctx context.Context, m *s_d.Organization) (*s_d.Organization, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.StaffRole {
		c.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	// Add defaults.
	m.ID = primitive.NewObjectID()
	m.CreatedByUserID = userID
	m.CreatedByUserName = userName
	m.CreatedAt = time.Now()
	m.ModifiedByUserID = userID
	m.ModifiedByUserName = userName
	m.ModifiedAt = time.Now()

	// Save to our database.
	err := c.OrganizationStorer.Create(ctx, m)
	if err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	return m, nil
}
