package controller

import (
	"context"

	"golang.org/x/exp/slog"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (c *UserControllerImpl) ListByFilter(ctx context.Context, f *user_s.UserListFilter) (*user_s.UserListResult, error) {
	// Extract from our session the following data.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole != user_s.UserRoleRoot {
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("Status", f.Status),
		slog.String("SearchText", f.SearchText),
		slog.Time("CreatedAtGTE", f.CreatedAtGTE),
		slog.Bool("ExcludeArchived", f.ExcludeArchived))

	// Filtering the database.
	m, err := c.UserStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}

func (c *UserControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *user_s.UserListFilter) ([]*user_s.UserAsSelectOption, error) {
	// Extract from our session the following data.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole != user_s.UserRoleRoot {
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("Role", f.Role))

	// Filtering the database.
	m, err := c.UserStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
