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
	SubmissionPendingState                        = 1
	SubmissionActiveState                         = 2
	SubmissionErrorState                          = 3
	SubmissionInactiveState                       = 4
	PreScreeningServiceType                       = 1
	PedigreeServiceType                           = 2
	CPSCapsuleYougradeServiceType                 = 3
	PoorFinding                                   = 1
	FairFinding                                   = 2
	GoodFinding                                   = 3
	VeryGoodFinding                               = 4
	FineFinding                                   = 5
	VeryFineFinding                               = 6
	NearMintFinding                               = 7
	YesItShowsSignsOfTamperingOrRestoration       = 1
	NoItDoesNotShowsSignsOfTamperingOrRestoration = 2
	LetterGradeScale                              = 1
	NumberGradeScale                              = 2
	CPSPercentageGradingScale                     = 3
)

type Submission struct {
	ID                                 primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt                          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt                         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ServiceType                        int8               `bson:"service_type" json:"service_type"`
	State                              int8               `bson:"state" json:"state"`
	SubmissionDate                     time.Time          `bson:"submission_date" json:"submission_date"`
	Item                               string             `bson:"item" json:"item"` // Created by system.
	SeriesTitle                        string             `bson:"series_title" json:"series_title"`
	IssueVol                           string             `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string             `bson:"issue_no" json:"issue_no"`
	IssueCoverDate                     string             `bson:"issue_cover_date" json:"issue_cover_date"`
	PublisherName                      string             `bson:"publisher_name" json:"publisher_name"`
	PublisherNameOther                 string             `bson:"publisher_name_other" json:"publisher_name_other"`
	SpecialNotesLine1                  string             `bson:"special_notes_line_1" json:"special_notes_line_1"`
	SpecialNotesLine2                  string             `bson:"special_notes_line_2" json:"special_notes_line_2"`
	SpecialNotesLine3                  string             `bson:"special_notes_line_3" json:"special_notes_line_3"`
	SpecialNotesLine4                  string             `bson:"special_notes_line_4" json:"special_notes_line_4"`
	SpecialNotesLine5                  string             `bson:"special_notes_line_5" json:"special_notes_line_5"`
	GradingNotesLine1                  string             `bson:"grading_notes_line_1" json:"grading_notes_line_1"`
	GradingNotesLine2                  string             `bson:"grading_notes_line_2" json:"grading_notes_line_2"`
	GradingNotesLine3                  string             `bson:"grading_notes_line_3" json:"grading_notes_line_3"`
	GradingNotesLine4                  string             `bson:"grading_notes_line_4" json:"grading_notes_line_4"`
	GradingNotesLine5                  string             `bson:"grading_notes_line_5" json:"grading_notes_line_5"`
	CreasesFinding                     string             `bson:"creases_finding" json:"creases_finding"`
	TearsFinding                       string             `bson:"tears_finding" json:"tears_finding"`
	MissingPartsFinding                string             `bson:"missing_parts_finding" json:"missing_parts_finding"`
	StainsFinding                      string             `bson:"stains_finding" json:"stains_finding"`
	DistortionFinding                  string             `bson:"distortion_finding" json:"distortion_finding"`
	PaperQualityFinding                string             `bson:"paper_quality_finding" json:"paper_quality_finding"`
	SpineFinding                       string             `bson:"spine_finding" json:"spine_finding"`
	CoverFinding                       string             `bson:"cover_finding" json:"cover_finding"`
	ShowsSignsOfTamperingOrRestoration int8               `bson:"shows_signs_of_tampering_or_restoration" json:"shows_signs_of_tampering_or_restoration"`
	GradingScale                       int8               `bson:"grading_scale" json:"grading_scale"`
	OverallLetterGrade                 string             `bson:"overall_letter_grade" json:"overall_letter_grade"`
	OverallNumberGrade                 float64            `bson:"overall_number_grade" json:"overall_number_grade"`
	CpsPercentageGrade                 float64            `bson:"cps_percentage_grade" json:"cps_percentage_grade"`
	UserID                             primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserFirstName                      string             `bson:"user_first_name" json:"user_first_name"`
	UserLastName                       string             `bson:"user_last_name" json:"user_last_name"`
	UserCompanyName                    string             `bson:"user_company_name" json:"user_company_name"`
	UserSignature                      string             `bson:"user_signature" json:"user_signature"`
	InspectorSignature                 string             `bson:"inspector_signature" json:"inspector_signature"`
	InspectorDate                      time.Time          `bson:"inspector_date" json:"inspector_date"`
	InspectorFirstName                 string             `bson:"inspector_first_name" json:"inspector_first_name"`
	InspectorLastName                  string             `bson:"inspector_last_name" json:"inspector_last_name"`
	InspectorCompany                   string             `bson:"inspector_company_name" json:"inspector_company_name"`
	SecondInspectorSignature           string             `bson:"second_inspector_signature" json:"second_inspector_signature"`
	SecondInspectorFirstName           string             `bson:"second_inspector_first_name" json:"second_inspector_first_name"`
	SecondInspectorLastName            string             `bson:"second_inspector_last_name" json:"second_inspector_last_name"`
	SecondInspectorCompany             string             `bson:"second_inspector_company" json:"second_inspector_company"`
	SecondInspectorDate                time.Time          `bson:"second_inspector_date" json:"second_inspector_date"`
	ThirdInspectorSignature            string             `bson:"third_inspector_signature" json:"third_inspector_signature"`
	ThirdInspectorFirstName            string             `bson:"third_inspector_first_name" json:"third_inspector_first_name"`
	ThirdInspectorLastName             string             `bson:"third_inspector_last_name" json:"third_inspector_last_name"`
	ThirdInspectorCompany              string             `bson:"third_inspector_company" json:"third_inspector_company"`
	ThirdInspectorDate                 time.Time          `bson:"third_inspector_date" json:"third_inspector_date"`
	Filename                           string             `bson:"filename" json:"filename"`
	FileUploadS3ObjectKey              string             `bson:"file_upload_s3_key" json:"file_upload_s3_object_key"`
	FileUploadDownloadableFileURL      string
}

type SubmissionListFilter struct {
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

type SubmissionListResult struct {
	Results []*Submission `json:"results"`
}

// SubmissionStorer Interface for submission.
type SubmissionStorer interface {
	Create(ctx context.Context, m *Submission) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Submission, error)
	UpdateByID(ctx context.Context, m *Submission) error
	ListByFilter(ctx context.Context, m *SubmissionListFilter) (*SubmissionListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
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
