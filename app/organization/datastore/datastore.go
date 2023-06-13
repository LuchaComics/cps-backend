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
	OrganizationPendingStatus  = 1
	OrganizationActiveStatus   = 2
	OrganizationErrorStatus    = 3
	OrganizationArchivedStatus = 4
	RetailerType               = 1
)

type Organization struct {
	ID                 primitive.ObjectID     `bson:"_id" json:"id"`
	ModifiedAt         time.Time              `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName string                 `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID   primitive.ObjectID     `bson:"modified_by_user_id" json:"modified_by_user_id"`
	Type               int8                   `bson:"type" json:"type"`
	Status             int8                   `bson:"status" json:"status"`
	Name               string                 `bson:"name" json:"name"` // Created by system.
	CreatedAt          time.Time              `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName  string                 `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID    primitive.ObjectID     `bson:"created_by_user_id" json:"created_by_user_id"`
	Comments           []*OrganizationComment `bson:"comments" json:"comments"`
}

type OrganizationComment struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID   primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserID  primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByName    string             `bson:"created_by_name" json:"created_by_name"`
	ModifiedAt       time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserID primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByName   string             `bson:"modified_by_name" json:"modified_by_name"`
	Content          string             `bson:"content" json:"content"`
}

type OrganizationListFilter struct {
	PageSize        int64
	LastID          string
	SortField       string
	UserID          primitive.ObjectID
	UserRole        int8
	ExcludeArchived bool

	// SortOrder string   `json:"sort_order"`
	// SortField string   `json:"sort_field"`
	// Offset    uint64   `json:"offset"`
	// Limit     uint64   `json:"limit"`
	// Statuss    []int8   `json:"statuss"`
	// UUIDs     []string `json:"uuids"`
}

type OrganizationListResult struct {
	Results []*Organization `json:"results"`
}

type OrganizationAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// OrganizationStorer Interface for organization.
type OrganizationStorer interface {
	Create(ctx context.Context, m *Organization) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Organization, error)
	UpdateByID(ctx context.Context, m *Organization) error
	ListByFilter(ctx context.Context, m *OrganizationListFilter) (*OrganizationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *OrganizationListFilter) ([]*OrganizationAsSelectOption, error)
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
