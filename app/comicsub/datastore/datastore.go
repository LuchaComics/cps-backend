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
	StatusPending                                 = 1
	StatusActive                                  = 2
	StatusError                                   = 3
	StatusArchived                                = 100
	ServiceTypePreScreening                       = 1
	ServiceTypePedigree                           = 2
	ServiceTypeCPSCapsule                         = 3
	ServiceTypeCPSCapsuleIndieMintGem             = 4
	ServiceTypeCPSCapsuleSignatureCollection      = 5
	FindingPoor                                   = 1
	FindingFair                                   = 2
	FindingGood                                   = 3
	FindingVeryGood                               = 4
	FindingFine                                   = 5
	FindingVeryFine                               = 6
	FindingNearMint                               = 7
	YesItShowsSignsOfTamperingOrRestoration       = 1
	NoItDoesNotShowsSignsOfTamperingOrRestoration = 2
	GradingScaleLetter                            = 1
	GradingScaleNumber                            = 2
	GradingScaleCPSPercentage                     = 3
	CollectibleTypeGeneric                        = 1
)

type ComicSubmission struct {
	ID                                 primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID                     primitive.ObjectID `bson:"organization_id,omitempty" json:"organization_id,omitempty"`
	OrganizationName                   string             `bson:"organization_name" json:"organization_name"`
	CPSRN                              string             `bson:"cpsrn" json:"cpsrn"`
	CreatedAt                          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserID                    primitive.ObjectID `bson:"created_by_user_id,omitempty" json:"created_by_user_id,omitempty"`
	CreatedByUserRole                  int8               `bson:"created_by_user_role" json:"created_by_user_role"`
	ModifiedAt                         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserID                   primitive.ObjectID `bson:"modified_by_user_id,omitempty" json:"modified_by_user_id,omitempty"`
	ModifiedByUserRole                 int8               `bson:"modified_by_user_role" json:"modified_by_user_role"`
	ServiceType                        int8               `bson:"service_type" json:"service_type"`
	Status                             int8               `bson:"status" json:"status"`
	SubmissionDate                     time.Time          `bson:"submission_date" json:"submission_date"`
	Item                               string             `bson:"item" json:"item"` // Created by system.
	SeriesTitle                        string             `bson:"series_title" json:"series_title"`
	IssueVol                           string             `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string             `bson:"issue_no" json:"issue_no"`
	IssueCoverYear                     int64              `bson:"issue_cover_year" json:"issue_cover_year"`
	IssueCoverMonth                    int8               `bson:"issue_cover_month" json:"issue_cover_month"`
	PublisherName                      int8               `bson:"publisher_name" json:"publisher_name"`
	PublisherNameOther                 string             `bson:"publisher_name_other" json:"publisher_name_other"`
	SpecialNotes                       string             `bson:"special_notes" json:"special_notes"`
	GradingNotes                       string             `bson:"grading_notes" json:"grading_notes"`
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
	IsOverallLetterGradeNearMintPlus   bool               `bson:"is_overall_letter_grade_near_mint_plus" json:"is_overall_letter_grade_near_mint_plus"`
	IsCpsIndieMintGem                  bool               `bson:"is_cps_indie_mint_gem" json:"is_cps_indie_mint_gem"`
	UserID                             primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"` // This is the customer this submission belongs to.
	User                               *SubmissionUser    `bson:"user" json:"user"`                           // This is the customer this submission belongs to.
	UserFirstName                      string             `bson:"user_first_name" json:"user_first_name"`
	UserLastName                       string             `bson:"user_last_name" json:"user_last_name"`
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
	Comments                           []*SubmissionComment   `bson:"comments" json:"comments,omitempty"`
	CollectibleType                    int8                   `bson:"collectible_type" json:"collectible_type"`
	Signatures                         []*SubmissionSignature `bson:"signatures" json:"signatures,omitempty"`
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

type SubmissionSignature struct {
	Role string `bson:"role" json:"role"`
	Name string `bson:"name" json:"name"`
}

type ComicSubmissionListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	OrganizationID    primitive.ObjectID
	UserID            primitive.ObjectID
	UserEmail         string
	CreatedByUserRole int8
	ExcludeArchived   bool
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
	Status                    int8               `bson:"status" json:"status"`
}

type ComicSubmissionListResult struct {
	Results     []*ComicSubmission `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type ComicSubmissionAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// ComicSubmissionStorer Interface for submission.
type ComicSubmissionStorer interface {
	Create(ctx context.Context, m *ComicSubmission) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*ComicSubmission, error)
	GetByCPSRN(ctx context.Context, cpsrn string) (*ComicSubmission, error)
	UpdateByID(ctx context.Context, m *ComicSubmission) error
	ListByFilter(ctx context.Context, f *ComicSubmissionListFilter) (*ComicSubmissionListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *ComicSubmissionListFilter) ([]*ComicSubmissionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CountAll(ctx context.Context) (int64, error)
	CountByFilter(ctx context.Context, f *ComicSubmissionListFilter) (int64, error)
	// //TODO: Add more...
}

type ComicSubmissionStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) ComicSubmissionStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("comic_submissions")

	s := &ComicSubmissionStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
