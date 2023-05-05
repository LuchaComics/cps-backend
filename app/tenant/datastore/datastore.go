package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
)

const (
	TenantActiveState       = 1
	TenantInactiveState     = 2
	TenantAdministratorRole = 1
	TenantTenantRole        = 2
)

type Tenant struct {
	ID       primitive.ObjectID `bson:"_id"`
	TenantID string             `bson:"tenant_id"`
	Name     string             `bson:"name"`
}

type TenantFilter struct {
	SortOrder string   `json:"sort_order"`
	SortField string   `json:"sort_field"`
	Offset    uint64   `json:"offset"`
	Limit     uint64   `json:"limit"`
	States    []int8   `json:"states"`
	UUIDs     []string `json:"uuids"`
}

// TenantStorer Interface for tenant.
type TenantStorer interface {
	Create(ctx context.Context, m *Tenant) error
	GetByTenantID(ctx context.Context, tenantID string) (*Tenant, error)
	GetByName(ctx context.Context, email string) (*Tenant, error)
	CheckIfExistsByName(ctx context.Context, name string) (bool, error)
	UpdateByTenantID(ctx context.Context, m *Tenant) error
	// //TODO: Add more...
}

type TenantStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) TenantStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("tenants")

	s := &TenantStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
