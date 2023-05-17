package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	s_d "github.com/LuchaComics/cps-backend/app/submission/datastore"
	u_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
)

func (c *SubmissionControllerImpl) Create(ctx context.Context, m *s_d.Submission) error {
	// Modify the submission based on role.
	userRole, ok := ctx.Value(constants.SessionUserRole).(int8)
	if ok {
		switch userRole {
		case u_d.RetailerRole:
			// Override state.
			m.State = s_d.SubmissionPendingState

			// Auto-assign the user-if
			m.UserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
			m.UserFirstName = ctx.Value(constants.SessionUserFirstName).(string)
			m.UserLastName = ctx.Value(constants.SessionUserLastName).(string)
			m.UserCompanyName = ctx.Value(constants.SessionUserCompanyName).(string)
			m.ServiceType = s_d.PreScreeningServiceType
		case u_d.StaffRole:
			m.State = s_d.SubmissionActiveState
		default:
			m.State = s_d.SubmissionErrorState
		}
	}

	// Add defaults.
	m.ID = primitive.NewObjectID()
	m.CreatedTime = time.Now()
	m.ModifiedTime = time.Now()
	m.SubmissionDate = time.Now()

	// Save to our database.
	err := c.SubmissionStorer.Create(ctx, m)
	if err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return err
	}

	r := &pdfbuilder.CBFFBuilderRequestDTO{
		ID:                                 m.ID,
		Filename:                           fmt.Sprintf("%v.pdf", m.ID.Hex()),
		SubmissionDate:                     time.Now(),
		SeriesTitle:                        m.SeriesTitle,
		IssueVol:                           m.IssueVol,
		IssueNo:                            m.IssueNo,
		IssueCoverDate:                     m.IssueCoverDate,
		PublisherName:                      m.PublisherName,
		SpecialNotesLine1:                  m.SpecialNotesLine1,
		SpecialNotesLine2:                  m.SpecialNotesLine2,
		SpecialNotesLine3:                  m.SpecialNotesLine3,
		SpecialNotesLine4:                  m.SpecialNotesLine4,
		SpecialNotesLine5:                  m.SpecialNotesLine5,
		GradingNotesLine1:                  m.GradingNotesLine1,
		GradingNotesLine2:                  m.GradingNotesLine2,
		GradingNotesLine3:                  m.GradingNotesLine3,
		GradingNotesLine4:                  m.GradingNotesLine4,
		GradingNotesLine5:                  m.GradingNotesLine5,
		CreasesFinding:                     m.CreasesFinding,
		TearsFinding:                       m.TearsFinding,
		MissingPartsFinding:                m.MissingPartsFinding,
		StainsFinding:                      m.StainsFinding,
		DistortionFinding:                  m.DistortionFinding,
		PaperQualityFinding:                m.PaperQualityFinding,
		SpineFinding:                       m.SpineFinding,
		CoverFinding:                       m.CoverFinding,
		ShowsSignsOfTamperingOrRestoration: m.ShowsSignsOfTamperingOrRestoration,
		OverallLetterGrade:                 m.OverallLetterGrade,
		UserFirstName:                      m.UserFirstName,
		UserLastName:                       m.UserLastName,
		UserCompanyName:                    m.UserCompanyName,
	}
	res, err := c.CBFFBuilder.GeneratePDF(r)
	log.Println("===--->", res, err, "<---===") //TODO: IMPL SAVING TO S3.

	return err
}
