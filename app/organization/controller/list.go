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

func (c *OrganizationControllerImpl) ListByFilter(ctx context.Context, f *domain.OrganizationListFilter) (*domain.OrganizationListResult, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		c.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	c.Logger.Debug("fetching organizations now...", slog.Any("userID", userID))

	m, err := c.OrganizationStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched organizations", slog.Any("m", m))
	return m, err
}

func (c *OrganizationControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.OrganizationListFilter) ([]*domain.OrganizationAsSelectOption, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		c.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	c.Logger.Debug("fetching organizations now...", slog.Any("userID", userID))

	m, err := c.OrganizationStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched organizations", slog.Any("m", m))
	return m, err
}
