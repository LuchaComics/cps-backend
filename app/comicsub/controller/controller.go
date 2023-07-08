package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	mg "github.com/LuchaComics/cps-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	s3_storage "github.com/LuchaComics/cps-backend/adapter/storage/s3"
	submission_s "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	organization_s "github.com/LuchaComics/cps-backend/app/organization/datastore"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/cpsrn"
	"github.com/LuchaComics/cps-backend/provider/kmutex"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// ComicSubmissionController Interface for submission business logic controller.
type ComicSubmissionController interface {
	Create(ctx context.Context, req *ComicSubmissionCreateRequestIDO) (*submission_s.ComicSubmission, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*submission_s.ComicSubmission, error)
	GetByCPSRN(ctx context.Context, cpsrn string) (*submission_s.ComicSubmission, error)
	UpdateByID(ctx context.Context, req *ComicSubmissionUpdateRequestIDO) (*submission_s.ComicSubmission, error)
	ListByFilter(ctx context.Context, f *submission_s.ComicSubmissionListFilter) (*submission_s.ComicSubmissionListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *submission_s.ComicSubmissionListFilter) ([]*submission_s.ComicSubmissionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*submission_s.ComicSubmission, error)
	SetUser(ctx context.Context, submissionID primitive.ObjectID, userID primitive.ObjectID) (*submission_s.ComicSubmission, error)
	CreateComment(ctx context.Context, submissionID primitive.ObjectID, content string) (*submission_s.ComicSubmission, error)
}

type ComicSubmissionControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	S3                    s3_storage.S3Storager
	Password              password.Provider
	CPSRN                 cpsrn.Provider
	CBFFBuilder           pdfbuilder.CBFFBuilder
	PCBuilder             pdfbuilder.PCBuilder
	CCIMGBuilder          pdfbuilder.CCIMGBuilder
	CCSCBuilder           pdfbuilder.CCSCBuilder
	CCBuilder             pdfbuilder.CCBuilder
	CCUGBuilder           pdfbuilder.CCUGBuilder
	Emailer               mg.Emailer
	Kmutex                kmutex.Provider
	UserStorer            user_s.UserStorer
	ComicSubmissionStorer submission_s.ComicSubmissionStorer
	OrganizationStorer    organization_s.OrganizationStorer
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
	pcb pdfbuilder.PCBuilder,
	ccimg pdfbuilder.CCIMGBuilder,
	ccsc pdfbuilder.CCSCBuilder,
	cc pdfbuilder.CCBuilder,
	ccug pdfbuilder.CCUGBuilder,
	emailer mg.Emailer,
	usr_storer user_s.UserStorer,
	sub_storer submission_s.ComicSubmissionStorer,
	org_storer organization_s.OrganizationStorer,
) ComicSubmissionController {
	loggerp.Debug("submission controller initialization started...")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // FOR TESTING PURPOSES ONLY.
	// text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis, unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam eaque ipsa, quae ab illo inventore veritatis et quasi architecto beatae`
	// r := &pdfbuilder.CBFFBuilderRequestDTO{
	// 	CPSRN:                              "788346-26649-1-1000",
	// 	ComicSubmissionDate:                     time.Now(),
	// 	SeriesTitle:                        "Winter World",
	// 	IssueVol:                           "Vol 1",
	// 	IssueNo:                            "#1",
	// 	IssueCoverYear:                     2023,
	// 	IssueCoverMonth:                    1,
	// 	PublisherName:                      "Some publisher",
	// 	SpecialNotes:                       text,
	// 	GradingNotes:                       text,
	// 	CreasesFinding:                     "VF",
	// 	TearsFinding:                       "FN",
	// 	MissingPartsFinding:                "PR",
	// 	StainsFinding:                      "NM",
	// 	DistortionFinding:                  "NM",
	// 	PaperQualityFinding:                "VF",
	// 	SpineFinding:                       "FN",
	// 	CoverFinding:                       "VG",
	// 	GradingScale:                       1,
	// 	ShowsSignsOfTamperingOrRestoration: true,
	// 	OverallLetterGrade:                 "NM",
	// 	IsOverallLetterGradeNearMintPlus:   true,
	// 	UserFirstName:                      "Bartlomiej",
	// 	UserLastName:                       "Miks",
	// 	UserOrganizationName:               "Mika Software Corporation",
	// }
	// res, err := cbffb.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // // FOR TESTING PURPOSES ONLY.
	// text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis, unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam eaque ipsa, quae ab illo inventore veritatis et quasi architecto beatae`
	// r := &pdfbuilder.PCBuilderRequestDTO{
	// 	CPSRN:                              "788346-26649-1-1000",
	// 	SeriesTitle:                        "Winter World",
	// 	IssueVol:                           "Vol 1",
	// 	IssueNo:                            "#1",
	// 	IssueCoverYear:                     2023,
	// 	IssueCoverMonth:                    1,
	// 	PublisherName:                      "Some publisher",
	// 	SpecialDetails:                     1, // 2=Regular Edition
	// 	SpecialDetailsOther:                "this is a test, lalalala",
	// 	SpecialNotes:                       text,
	// 	GradingNotes:                       text,
	// 	CreasesFinding:                     "PR",
	// 	TearsFinding:                       "FN",
	// 	MissingPartsFinding:                "PR",
	// 	StainsFinding:                      "NM",
	// 	DistortionFinding:                  "NM",
	// 	PaperQualityFinding:                "VF",
	// 	SpineFinding:                       "FN",
	// 	CoverFinding:                       "VG",
	// 	GradingScale:                       1,
	// 	ShowsSignsOfTamperingOrRestoration: true,
	// 	OverallLetterGrade:                 "NM",
	// 	IsOverallLetterGradeNearMintPlus:   true,
	// 	UserFirstName:                      "Bartlomiej",
	// 	UserLastName:                       "Miks",
	// 	UserOrganizationName:               "Mika Software Corporation",
	// }
	// res, err := pcb.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CCIMGBuilderRequestDTO{
	// 	CPSRN:           "788346-26649-1-1000",
	// 	SeriesTitle:     "Winter World",
	// 	IssueVol:        "Vol 1",
	// 	IssueNo:         "#1",
	// 	IssueCoverYear:  2023,
	// 	IssueCoverMonth: 1,
	// 	PublisherName:   "Some publisher",
	// 	SpecialDetails:  2, // 2=Regular Edition
	//
	// }
	// res, err := ccimg.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CCSCBuilderRequestDTO{
	// 	CPSRN:                              "788346-26649-1-1000",
	// 	SeriesTitle:                        "Winter World",
	// 	IssueVol:                           "Vol 1",
	// 	IssueNo:                            "#1",
	// 	IssueCoverYear:                     2023,
	// 	IssueCoverMonth:                    1,
	// 	PublisherName:                      "Some publisher",
	// 	SpecialDetails:                     2, // 2=Regular Edition
	// 	GradingScale:                       1,
	// 	ShowsSignsOfTamperingOrRestoration: true,
	// 	OverallLetterGrade:                 "NM",
	// 	IsOverallLetterGradeNearMintPlus:   true,
	// }
	// res, err := ccsc.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CCBuilderRequestDTO{
	// 	CPSRN:                            "788346-26649-1-1000",
	// 	SeriesTitle:                      "Winter World",
	// 	IssueVol:                         "Vol 1",
	// 	IssueNo:                          "#1",
	// 	IssueCoverYear:                   2023,
	// 	IssueCoverMonth:                  1,
	// 	PublisherName:                    "Some publisher",
	// 	SpecialDetails:                   2, // 2=Regular Edition
	// 	GradingScale:                     1,
	// 	OverallLetterGrade:               "vf",
	// 	IsOverallLetterGradeNearMintPlus: false,
	// 	OverallNumberGrade:               10,
	// 	CpsPercentageGrade:               100,
	// }
	// res, err := cc.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//
	//------------------------------------------------------------------------//

	// // // FOR TESTING PURPOSES ONLY.
	// r := &pdfbuilder.CCUGBuilderRequestDTO{
	// 	CPSRN:                            "788346-26649-1-1000",
	// 	SeriesTitle:                      "Winter World",
	// 	IssueVol:                         "Vol 1",
	// 	IssueNo:                          "#1",
	// 	IssueCoverYear:                   2023,
	// 	IssueCoverMonth:                  1,
	// 	PublisherName:                    "Some publisher",
	// 	SpecialDetails:                   2, // 2=Regular Edition
	// 	GradingScale:                     3, // 1=Letter 2=Number 3=CPS
	// 	OverallLetterGrade:               "vf",
	// 	IsOverallLetterGradeNearMintPlus: false,
	// 	OverallNumberGrade:               7,
	// 	CpsPercentageGrade:               100,
	// }
	// res, err := ccug.GeneratePDF(r)
	// log.Println("===--->", res, err, "<---===")

	// ------------------------------------------------------------------------//
	// ------------------------------------------------------------------------//
	// ------------------------------------------------------------------------//

	s := &ComicSubmissionControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		S3:                    s3,
		Password:              passwordp,
		Kmutex:                kmux,
		CPSRN:                 cpsrnP,
		CBFFBuilder:           cbffb,
		PCBuilder:             pcb,
		CCIMGBuilder:          ccimg,
		CCSCBuilder:           ccsc,
		CCBuilder:             cc,
		CCUGBuilder:           ccug,
		Emailer:               emailer,
		UserStorer:            usr_storer,
		ComicSubmissionStorer: sub_storer,
		OrganizationStorer:    org_storer,
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
		Status:                    u.Status,
		Role:                      u.Role,
	}
}
