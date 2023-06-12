package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	user_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (c *SubmissionControllerImpl) ListByFilter(ctx context.Context, f *domain.SubmissionListFilter) (*domain.SubmissionListResult, error) {
	// Extract from our session the following data.
	organizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on tenancy if the user is not a system administrator.
	if userRole != user_d.StaffRole {
		f.OrganizationID = organizationID
		c.Logger.Debug("applying security policy to filters",
			slog.Any("organization_id", organizationID),
			slog.Any("user_id", userID),
			slog.Any("user_role", userRole))
	}

	m, err := c.SubmissionStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
