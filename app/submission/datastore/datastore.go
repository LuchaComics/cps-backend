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
	SubmissionArchivedState                       = 100
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
	OrganizationID                     primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	CPSRN                              string             `bson:"cpsrn" json:"cpsrn"`
	CreatedAt                          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserID                    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByUserRole                  int8               `bson:"created_by_user_role" json:"created_by_user_role"`
	ModifiedAt                         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserID                   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByUserRole                 int8               `bson:"modified_by_user_role" json:"modified_by_user_role"`
	ServiceType                        int8               `bson:"service_type" json:"service_type"`
	State                              int8               `bson:"state" json:"state"`
	SubmissionDate                     time.Time          `bson:"submission_date" json:"submission_date"`
	Item                               string             `bson:"item" json:"item"` // Created by system.
	SeriesTitle                        string             `bson:"series_title" json:"series_title"`
	IssueVol                           string             `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string             `bson:"issue_no" json:"issue_no"`
	IssueCoverYear                     int64              `bson:"issue_cover_year" json:"issue_cover_year"`
	IssueCoverMonth                    int8               `bson:"issue_cover_month" json:"issue_cover_month"`
	PublisherName                      int8               `bson:"publisher_name" json:"publisher_name"`
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
	UserID                             primitive.ObjectID `bson:"user_id" json:"user_id"` // This is the customer this submission belongs to.
	User                               *SubmissionUser    `bson:"user" json:"user"`       // This is the customer this submission belongs to.
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
	Comments                           []*SubmissionComment `bson:"comments" json:"comments,omitempty"`
}

type SubmissionComment struct {
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

type SubmissionListFilter struct {
	PageSize        int64
	LastID          string
	SortField       string
	OrganizationID  primitive.ObjectID
	UserID          primitive.ObjectID
	UserRole        int8
	ExcludeArchived bool
}

type SubmissionUser struct {
	ID                        primitive.ObjectID `bson:"_id" json:"_id"`
	OrganizationID            primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	FirstName                 string             `bson:"first_name" json:"first_name"`
	LastName                  string             `bson:"last_name" json:"last_name"`
	Name                      string             `bson:"name" json:"name"`
	LexicalName               string             `bson:"lexical_name" json:"lexical_name"`
	Email                     string             `bson:"email" json:"email"`
	Phone                     string             `bson:"phone" json:"phone,omitempty"`
	Country                   string             `bson:"country" json:"country,omitempty"`
	Region                    string             `bson:"region" json:"region,omitempty"`
	City                      string             `bson:"city" json:"city,omitempty"`
	PostalCode                string             `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1              string             `bson:"address_line_1" json:"address_line_1,omitempty"`
	AddressLine2              string             `bson:"address_line_2" json:"address_line_2,omitempty"`
	HowDidYouHearAboutUs      int8               `bson:"how_did_you_hear_about_us" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `bson:"how_did_you_hear_about_us_other" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreePromotionsEmail      bool               `bson:"agree_promotions_email" json:"agree_promotions_email,omitempty"`
	CreatedAt                 time.Time          `bson:"created_at" json:"created_at,omitempty"`
	ModifiedAt                time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	Role                      int8               `bson:"role" json:"role"`
	State                     int8               `bson:"state" json:"state"`
}

type SubmissionListResult struct {
	Results []*Submission `json:"results"`
}

// SubmissionStorer Interface for submission.
type SubmissionStorer interface {
	Create(ctx context.Context, m *Submission) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Submission, error)
	GetByCPSRN(ctx context.Context, cpsrn string) (*Submission, error)
	UpdateByID(ctx context.Context, m *Submission) error
	ListByFilter(ctx context.Context, f *SubmissionListFilter) (*SubmissionListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CountAll(ctx context.Context) (int64, error)
	CountByFilter(ctx context.Context, f *SubmissionListFilter) (int64, error)
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
