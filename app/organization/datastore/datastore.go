package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
)

const (
	OrganizationPendingState  = 1
	OrganizationActiveState   = 2
	OrganizationErrorState    = 3
	OrganizationInactiveState = 4
	RetailerType              = 1
)

type Organization struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt       time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt      time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	Type            int8               `bson:"type" json:"type"`
	State           int8               `bson:"state" json:"state"`
	Name            string             `bson:"name" json:"name"` // Created by system.
	CreatedByUserID primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
}

type OrganizationListFilter struct {
	PageSize  int64
	LastID    string
	SortField string
	UserID    primitive.ObjectID
	UserRole  int8

	// SortOrder string   `json:"sort_order"`
	// SortField string   `json:"sort_field"`
	// Offset    uint64   `json:"offset"`
	// Limit     uint64   `json:"limit"`
	// States    []int8   `json:"states"`
	// UUIDs     []string `json:"uuids"`
}

type OrganizationListResult struct {
	Results []*Organization `json:"results"`
}

// OrganizationStorer Interface for organization.
type OrganizationStorer interface {
	Create(ctx context.Context, m *Organization) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Organization, error)
	UpdateByID(ctx context.Context, m *Organization) error
	ListByFilter(ctx context.Context, m *OrganizationListFilter) (*OrganizationListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type OrganizationStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) OrganizationStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("organizations")

	s := &OrganizationStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
