package controller

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	s_d "github.com/LuchaComics/cps-backend/app/submission/datastore"
	u_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
)

func (c *SubmissionControllerImpl) Create(ctx context.Context, m *s_d.Submission) (*s_d.Submission, error) {
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
		return nil, err
	}

	// The next following lines of code will create the PDF file gnerator
	// request to be submitted into our PDF file generator to generate the data.
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
		ShowsSignsOfTamperingOrRestoration: m.ShowsSignsOfTamperingOrRestoration == 1,
		GradingScale:                       m.GradingScale,
		OverallLetterGrade:                 m.OverallLetterGrade,
		OverallNumberGrade:                 m.OverallNumberGrade,
		CpsPercentageGrade:                 m.CpsPercentageGrade,
		UserFirstName:                      m.UserFirstName,
		UserLastName:                       m.UserLastName,
		UserCompanyName:                    m.UserCompanyName,
	}
	response, err := c.CBFFBuilder.GeneratePDF(r)
	if err != nil {
		c.Logger.Error("generate pdf error", slog.Any("error", err))
		return nil, err
	}
	if response == nil {
		c.Logger.Error("generate pdf error does not return a response")
		return nil, errors.New("no response from pdf generator")
	}

	// The next few lines will upload our PDF to our remote storage. Once the
	// file is saved remotely, we will have a connection to it through a "key"
	// unique reference to the uploaded file.
	path := fmt.Sprintf("uploads/%v", response.FileName)
	err = c.S3.UploadContent(ctx, path, response.Content)
	if err != nil {
		c.Logger.Error("s3 upload error", slog.Any("error", err))
		return nil, err
	}

	// The following will save the S3 key of our file upload into our record.
	m.FileUploadS3ObjectKey = path
	m.ModifiedTime = time.Now()

	if err := c.SubmissionStorer.UpdateByID(ctx, m); err != nil {
		c.Logger.Error("database update error", slog.Any("error", err))
		return nil, err
	}

	// The following will generate a pre-signed URL so user can download the file.
	downloadableURL, err := c.S3.GetDownloadablePresignedURL(ctx, m.FileUploadS3ObjectKey, time.Minute*15)
	if err != nil {
		c.Logger.Error("s3 presign error", slog.Any("error", err))
		return nil, err
	}
	m.FileUploadDownloadableFileURL = downloadableURL

	// Removing local file from the directory and don't do anything if we have errors.
	if err := os.Remove(response.FilePath); err != nil {
		c.Logger.Warn("removing local file error", slog.Any("error", err))
		// Just continue even if we get an error...
	}

	return m, nil
}
