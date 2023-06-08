package controller

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	mg "github.com/LuchaComics/cps-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	s3_storage "github.com/LuchaComics/cps-backend/adapter/storage/s3"
	organization_s "github.com/LuchaComics/cps-backend/app/organization/datastore"
	submission_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/cpsrn"
	"github.com/LuchaComics/cps-backend/provider/kmutex"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// SubmissionController Interface for submission business logic controller.
type SubmissionController interface {
	Create(ctx context.Context, m *submission_s.Submission) (*submission_s.Submission, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*submission_s.Submission, error)
	GetByCPSRN(ctx context.Context, cpsrn string) (*submission_s.Submission, error)
	UpdateByID(ctx context.Context, m *submission_s.Submission) (*submission_s.Submission, error)
	ListByFilter(ctx context.Context, f *submission_s.SubmissionListFilter) (*submission_s.SubmissionListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*submission_s.Submission, error)
	SetUser(ctx context.Context, submissionID primitive.ObjectID, userID primitive.ObjectID) (*submission_s.Submission, error)
	CreateComment(ctx context.Context, submissionID primitive.ObjectID, content string) (*submission_s.Submission, error)
}

type SubmissionControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	S3                 s3_storage.S3Storager
	Password           password.Provider
	CPSRN              cpsrn.Provider
	CBFFBuilder        pdfbuilder.CBFFBuilder
	Emailer            mg.Emailer
	Kmutex             kmutex.Provider
	UserStorer         user_s.UserStorer
	SubmissionStorer   submission_s.SubmissionStorer
	OrganizationStorer organization_s.OrganizationStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	kmux kmutex.Provider,
	cpsrnP cpsrn.Provider,
	cbffb pdfbuilder.CBFFBuilder,
	emailer mg.Emailer,
	usr_storer user_s.UserStorer,
	sub_storer submission_s.SubmissionStorer,
	org_storer organization_s.OrganizationStorer,
) SubmissionController {
	loggerp.Debug("submission controller initialization started...")

	// FOR TESTING PURPOSES ONLY.
	r := &pdfbuilder.CBFFBuilderRequestDTO{
		CPSRN:                              "788346-26649-1-1000",
		SubmissionDate:                     time.Now(),
		SeriesTitle:                        "Winter World",
		IssueVol:                           "Vol 1",
		IssueNo:                            "#1",
		IssueCoverYear:                     "2023",
		IssueCoverMonth:                    1,
		PublisherName:                      "Some publisher",
		SpecialNotesLine1:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		SpecialNotesLine2:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		SpecialNotesLine3:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		SpecialNotesLine4:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		SpecialNotesLine5:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		GradingNotesLine1:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		GradingNotesLine2:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		GradingNotesLine3:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		GradingNotesLine4:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		GradingNotesLine5:                  "XXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYYYYY",
		CreasesFinding:                     "VF",
		TearsFinding:                       "FN",
		MissingPartsFinding:                "PR",
		StainsFinding:                      "NM",
		DistortionFinding:                  "NM",
		PaperQualityFinding:                "VF",
		SpineFinding:                       "FN",
		CoverFinding:                       "VG",
		GradingScale:                       1,
		ShowsSignsOfTamperingOrRestoration: true,
		OverallLetterGrade:                 "VG",
		UserFirstName:                      "Bartlomiej",
		UserLastName:                       "Miks",
		UserCompanyName:                    "Mika Software Corporation",
	}
	res, err := cbffb.GeneratePDF(r)
	log.Println("===--->", res, err, "<---===")

	s := &SubmissionControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		S3:                 s3,
		Password:           passwordp,
		Kmutex:             kmux,
		CPSRN:              cpsrnP,
		CBFFBuilder:        cbffb,
		Emailer:            emailer,
		UserStorer:         usr_storer,
		SubmissionStorer:   sub_storer,
		OrganizationStorer: org_storer,
	}
	s.Logger.Debug("submission controller initialized")
	return s
}

// userToSubmissionUserCopy converts the full `User` record into a limited `SubmissionUser`.
func userToSubmissionUserCopy(u *user_s.User) *submission_s.SubmissionUser {
	if u == nil { // Defensive code.
		return nil
	}
	return &submission_s.SubmissionUser{
		ID:                        u.ID,
		OrganizationID:            u.OrganizationID,
		FirstName:                 u.FirstName,
		LastName:                  u.LastName,
		Name:                      u.Name,
		LexicalName:               u.LexicalName,
		Email:                     u.Email,
		Phone:                     u.Phone,
		Country:                   u.Country,
		Region:                    u.Region,
		City:                      u.City,
		PostalCode:                u.PostalCode,
		AddressLine1:              u.AddressLine1,
		AddressLine2:              u.AddressLine2,
		HowDidYouHearAboutUs:      u.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther: u.HowDidYouHearAboutUsOther,
		AgreePromotionsEmail:      u.AgreePromotionsEmail,
		CreatedAt:                 u.CreatedAt,
		ModifiedAt:                u.ModifiedAt,
		State:                     u.State,
		Role:                      u.Role,
	}
}
