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
	SubmissionActiveState         = 1
	SubmissionInactiveState       = 2
	PreScreeningServiceType       = 1
	PedigreeServiceType           = 2
	CPSCapsuleYougradeServiceType = 3
)

type Submission struct {
	ID                  primitive.ObjectID `bson:"_id"`
	UserID              string             `bson:"user_id" json:"user_id"`
	SubmissionID        string             `bson:"submission_id" json:"submission_id"`
	CreatedTime         time.Time          `bson:"created_time,omitempty" json:"created_time,omitempty"`
	ModifiedTime        time.Time          `bson:"modified_time,omitempty" json:"modified_time,omitempty"`
	ServiceType         int8               `bson:"service_type" json:"service_type"`
	Item                string             `bson:"item" json:"item"`
	Date                time.Time          `bson:"date" json:"date"`
	IssueTitle          string             `bson:"issue_title" json:"issue_title"`
	IssueVol            string             `bson:"issue_vol" json:"issue_vol"`
	IssueNo             string             `bson:"issue_no" json:"issue_no"`
	IssueDate           string             `bson:"issue_date" json:"issue_date"`
	IssueSpecialDetails string             `bson:"issue_special_details" json:"issue_special_details"`
}

type SubmissionFilter struct {
	SortOrder string   `json:"sort_order"`
	SortField string   `json:"sort_field"`
	Offset    uint64   `json:"offset"`
	Limit     uint64   `json:"limit"`
	States    []int8   `json:"states"`
	UUIDs     []string `json:"uuids"`
}

// SubmissionStorer Interface for submission.
type SubmissionStorer interface {
	Create(ctx context.Context, m *Submission) error
	GetBySubmissionID(ctx context.Context, submissionID string) (*Submission, error)
	UpdateBySubmissionID(ctx context.Context, m *Submission) error
	// //TODO: Add more...
}

type SubmissionStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) SubmissionStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("submissions")

	s := &SubmissionStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
