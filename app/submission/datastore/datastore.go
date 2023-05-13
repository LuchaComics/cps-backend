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
	SubmissionPendingState        = 1
	SubmissionActiveState         = 2
	SubmissionErrorState          = 3
	SubmissionInactiveState       = 4
	PreScreeningServiceType       = 1
	PedigreeServiceType           = 2
	CPSCapsuleYougradeServiceType = 3
	PoorFinding                   = 1
	FairFinding                   = 2
	GoodFinding                   = 3
	VeryGoodFinding               = 4
	FineFinding                   = 5
	VeryFineFinding               = 6
	NearMintFinding               = 7
)

type Submission struct {
	ID                       primitive.ObjectID `bson:"_id" json:"id"`
	SubmissionID             string             `bson:"submission_id" json:"submission_id"`
	CreatedTime              time.Time          `bson:"created_time,omitempty" json:"created_time,omitempty"`
	ModifiedTime             time.Time          `bson:"modified_time,omitempty" json:"modified_time,omitempty"`
	ServiceType              int8               `bson:"service_type" json:"service_type"`
	State                    int8               `bson:"state" json:"state"`
	Item                     string             `bson:"item" json:"item"` // Created by system.
	SeriesTitle              string             `bson:"series_title" json:"series_title"`
	IssueVol                 string             `bson:"issue_vol" json:"issue_vol"`
	IssueNo                  string             `bson:"issue_no" json:"issue_no"`
	IssueCoverDate           string             `bson:"issue_cover_date" json:"issue_cover_date"`
	IssueSpecialDetails      string             `bson:"issue_special_details" json:"issue_special_details"`
	CreasesFinding           string             `bson:"creases_finding" json:"creases_finding"`
	TearsFinding             string             `bson:"tears_finding" json:"tears_finding"`
	MissingPartsFinding      string             `bson:"missing_parts_finding" json:"missing_parts_finding"`
	StainsFinding            string             `bson:"stains_finding" json:"stains_finding"`
	DistortionFinding        string             `bson:"distortion_finding" json:"distortion_finding"`
	PaperQualityFinding      string             `bson:"paper_quality_finding" json:"paper_quality_finding"`
	SpineFinding             string             `bson:"spine_finding" json:"spine_finding"`
	CoverFinding             string             `bson:"cover_finding" json:"cover_finding"`
	OtherFinding             string             `bson:"other_finding" json:"other_finding"`
	OtherFindingText         string             `bson:"other_finding_text" json:"other_finding_text"`
	OverallLetterGrade       string             `bson:"overall_letter_grade" json:"overall_letter_grade"`
	UserID                   string             `bson:"user_id" json:"user_id"`
	UserFirstName            string             `bson:"user_first_name" json:"user_first_name"`
	UserLastName             string             `bson:"user_last_name" json:"user_last_name"`
	UserCompanyName          string             `bson:"user_company_name" json:"user_company_name"`
	UserSignature            string             `bson:"user_signature" json:"user_signature"`
	InspectorSignature       string             `bson:"inspector_signature" json:"inspector_signature"`
	InspectorDate            time.Time          `bson:"inspector_date" json:"inspector_date"`
	InspectorFirstName       string             `bson:"inspector_first_name" json:"inspector_first_name"`
	InspectorLastName        string             `bson:"inspector_last_name" json:"inspector_last_name"`
	InspectorCompany         string             `bson:"inspector_company_name" json:"inspector_company_name"`
	SecondInspectorSignature string             `bson:"second_inspector_signature" json:"second_inspector_signature"`
	SecondInspectorFirstName string             `bson:"second_inspector_first_name" json:"second_inspector_first_name"`
	SecondInspectorLastName  string             `bson:"second_inspector_last_name" json:"second_inspector_last_name"`
	SecondInspectorCompany   string             `bson:"second_inspector_company" json:"second_inspector_company"`
	SecondInspectorDate      time.Time          `bson:"second_inspector_date" json:"second_inspector_date"`
	ThirdInspectorSignature  string             `bson:"third_inspector_signature" json:"third_inspector_signature"`
	ThirdInspectorFirstName  string             `bson:"third_inspector_first_name" json:"third_inspector_first_name"`
	ThirdInspectorLastName   string             `bson:"third_inspector_last_name" json:"third_inspector_last_name"`
	ThirdInspectorCompany    string             `bson:"third_inspector_company" json:"third_inspector_company"`
	ThirdInspectorDate       time.Time          `bson:"third_inspector_date" json:"third_inspector_date"`
	// Note: Add company logo.
}

type SubmissionListFilter struct {
	PageSize  int64
	LastID    string
	SortField string
	UserID    primitive.ObjectID

	// SortOrder string   `json:"sort_order"`
	// SortField string   `json:"sort_field"`
	// Offset    uint64   `json:"offset"`
	// Limit     uint64   `json:"limit"`
	// States    []int8   `json:"states"`
	// UUIDs     []string `json:"uuids"`
}

type SubmissionListResult struct {
	Results []*Submission `json:"results"`
}

// SubmissionStorer Interface for submission.
type SubmissionStorer interface {
	Create(ctx context.Context, m *Submission) error
	GetByID(ctx context.Context, id string) (*Submission, error)
	UpdateByID(ctx context.Context, m *Submission) error
	ListByFilter(ctx context.Context, m *SubmissionListFilter) (*SubmissionListResult, error)
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
