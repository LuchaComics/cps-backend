package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	mg "github.com/LuchaComics/cps-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	s3_storage "github.com/LuchaComics/cps-backend/adapter/storage/s3"
	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	submission_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// SubmissionController Interface for submission business logic controller.
type SubmissionController interface {
	Create(ctx context.Context, m *domain.Submission) (*domain.Submission, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Submission, error)
	UpdateByID(ctx context.Context, m *domain.Submission) (*domain.Submission, error)
	ListByFilter(ctx context.Context, f *domain.SubmissionListFilter) (*domain.SubmissionListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type SubmissionControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	S3               s3_storage.S3Storager
	Password         password.Provider
	CBFFBuilder      pdfbuilder.CBFFBuilder
	Emailer          mg.Emailer
	SubmissionStorer submission_s.SubmissionStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	cbffb pdfbuilder.CBFFBuilder,
	emailer mg.Emailer,
	sub_storer submission_s.SubmissionStorer,
) SubmissionController {

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
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		S3:               s3,
		Password:         passwordp,
		CBFFBuilder:      cbffb,
		Emailer:          emailer,
		SubmissionStorer: sub_storer,
	}
	s.Logger.Debug("submission controller initialization started...")
	s.Logger.Debug("submission controller initialized")
	return s
}
