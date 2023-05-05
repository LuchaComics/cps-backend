package controller

import (
	"context"

	domain "github.com/LuchaComics/cps-backend/app/tenant/datastore"
)

func (c *TenantControllerImpl) GetTenantBySessionUUID(ctx context.Context, sessionUUID string) (*domain.Tenant, error) {
	panic("TODO: IMPLEMENT")
	return nil, nil
}
