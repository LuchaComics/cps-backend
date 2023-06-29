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
	StatusActive              = 1
	StatusError               = 2
	StatusArchived            = 3
	OwnershipTypeUser         = 1
	OwnershipTypeSubmission   = 2
	OwnershipTypeOrganization = 3
)

type Attachment struct {
	OrganizationID     primitive.ObjectID `bson:"organization_id,omitempty" json:"organization_id,omitempty"`
	OrganizationName   string             `bson:"organization_name" json:"organization_name"`
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	ModifiedAt         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	Name               string             `bson:"name" json:"name"`
	Description        string             `bson:"description" json:"description"`
	Filename           string             `bson:"filename" json:"filename"`
	ObjectKey          string             `bson:"object_key" json:"object_key"`
	ObjectURL          string             `bson:"object_url" json:"object_url"`
	OwnershipID        primitive.ObjectID `bson:"ownership_id" json:"ownership_id"`
	OwnershipType      int8               `bson:"ownership_type" json:"ownership_type"`
	Status             int8               `bson:"status" json:"status"`
}

type AttachmentListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID  primitive.ObjectID
	OwnershipID     primitive.ObjectID
	UserID          primitive.ObjectID
	UserRole        int8
	ExcludeArchived bool
}

type AttachmentListResult struct {
	Results     []*Attachment      `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type AttachmentAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// AttachmentStorer Interface for attachment.
type AttachmentStorer interface {
	Create(ctx context.Context, m *Attachment) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Attachment, error)
	UpdateByID(ctx context.Context, m *Attachment) error
	ListByFilter(ctx context.Context, m *AttachmentListFilter) (*AttachmentListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *AttachmentListFilter) ([]*AttachmentAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type AttachmentStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) AttachmentStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("attachments")

	s := &AttachmentStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
