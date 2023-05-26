package controller

import (
	"context"
	"time"

	domain "github.com/LuchaComics/cps-backend/app/organization/datastore"
	"golang.org/x/exp/slog"
)

func (c *OrganizationControllerImpl) UpdateByID(ctx context.Context, ns *domain.Organization) (*domain.Organization, error) {
	// Fetch the original organization.
	os, err := c.OrganizationStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, nil
	}

	// Modify our original organization.
	os.ModifiedAt = time.Now()
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
