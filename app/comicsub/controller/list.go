package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	user_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (c *ComicSubmissionControllerImpl) ListByFilter(ctx context.Context, f *domain.ComicSubmissionListFilter) (*domain.ComicSubmissionListResult, error) {
	// Extract from our session the following data.
	organizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on tenancy if the user is not a system administrator.
	if userRole != user_d.UserRoleRoot {
		f.OrganizationID = organizationID
		c.Logger.Debug("applying security policy to filters",
			slog.Any("organization_id", organizationID),
			slog.Any("user_id", userID),
			slog.Any("user_role", userRole))
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("OrganizationID", f.OrganizationID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("Status", f.Status),
		slog.Time("CreatedAtGTE", f.CreatedAtGTE),
		slog.String("SearchText", f.SearchText),
		slog.Bool("ExcludeArchived", f.ExcludeArchived))

	m, err := c.ComicSubmissionStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}

func (c *ComicSubmissionControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.ComicSubmissionListFilter) ([]*domain.ComicSubmissionAsSelectOption, error) {
	// Extract from our session the following data.
	organizationID := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on tenancy if the user is not a system administrator.
	if userRole != user_d.UserRoleRoot {
		f.OrganizationID = organizationID
		c.Logger.Debug("applying security policy to filters",
			slog.Any("organization_id", organizationID),
			slog.Any("user_id", userID),
			slog.Any("user_role", userRole))
	}

	m, err := c.ComicSubmissionStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list as select option by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
