package datastore

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	RootType                   = 1
	RetailerType               = 2
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
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID  primitive.ObjectID
	UserID          primitive.ObjectID
	UserRole        int8
	Status          int8
	ExcludeArchived bool
	SearchText      string
}

type OrganizationListResult struct {
	Results     []*Organization    `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
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

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &OrganizationStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
