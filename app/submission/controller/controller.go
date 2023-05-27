package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	mg "github.com/LuchaComics/cps-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	s3_storage "github.com/LuchaComics/cps-backend/adapter/storage/s3"
	organization_s "github.com/LuchaComics/cps-backend/app/organization/datastore"
	submission_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
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
	UpdateByID(ctx context.Context, m *submission_s.Submission) (*submission_s.Submission, error)
	ListByFilter(ctx context.Context, f *submission_s.SubmissionListFilter) (*submission_s.SubmissionListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
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
	sub_storer submission_s.SubmissionStorer,
	org_storer organization_s.OrganizationStorer,
) SubmissionController {
	loggerp.Debug("submission controller initialization started...")

	// // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CBFFBuilderRequestDTO{
	// 	ID:                                 primitive.NewObjectID(),
	// 	SubmissionDate:                     time.Now(),
	// 	SeriesTitle:                        "Winter World",
	// 	IssueVol:                           "Vol 1",
	// 	IssueNo:                            "#1",
	// 	IssueCoverDate:                     "16/05/2000",
	// 	PublisherName:                      "Some publisher",
	// 	SpecialNotesLine1:                  "XXXXXXXXXXXXXXXXX",
	// 	SpecialNotesLine2:                  "XXXXXXXXXXXXXXXXX",
	// 	SpecialNotesLine3:                  "XXXXXXXXXXXXXXXXX",
	// 	SpecialNotesLine4:                  "XXXXXXXXXXXXXXXXX",
	// 	SpecialNotesLine5:                  "XXXXXXXXXXXXXXXXX",
	// 	GradingNotesLine1:                  "XXXXXXXXXXXXXXXXX",
	// 	GradingNotesLine2:                  "XXXXXXXXXXXXXXXXX",
	// 	GradingNotesLine3:                  "XXXXXXXXXXXXXXXXX",
	// 	GradingNotesLine4:                  "XXXXXXXXXXXXXXXXX",
	// 	GradingNotesLine5:                  "XXXXXXXXXXXXXXXXX",
	// 	CreasesFinding:                     "VF",
	// 	TearsFinding:                       "FN",
	// 	MissingPartsFinding:                "PR",
	// 	StainsFinding:                      "NM",
	// 	DistortionFinding:                  "NM",
	// 	PaperQualityFinding:                "VF",
	// 	SpineFinding:                       "FN",
	// 	CoverFinding:                       "VG",
	// 	ShowsSignsOfTamperingOrRestoration: true,
	// 	OverallLetterGrade:                 "VG",
	// 	UserFirstName:                      "Bartlomiej",
	// 	UserLastName:                       "Miks",
	// 	UserCompanyName:                    "Mika Software Corporation",
	// }
	// res, err := cbffb.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

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
		SubmissionStorer:   sub_storer,
		OrganizationStorer: org_storer,
	}
	s.Logger.Debug("submission controller initialized")
	return s
}
