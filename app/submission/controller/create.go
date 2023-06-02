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
	// DEVELOPERS NOTE:
	// Every submission creation is dependent on the `role` of the logged in
	// user in our system so we need to extract it right away.
	userRole, ok := ctx.Value(constants.SessionUserRole).(int8)
	if !ok {
		c.Logger.Error("user role not extracted from session")
		return nil, fmt.Errorf("user role not extracted from session for submission id: %v", m.ID)
	}

	// DEVELOPERS NOTE:
	// Every submission needs to have a unique `CPS Registry Number` (CPRN)
	// generated. The following needs to happen to generate the unique CPRN:
	// 1. Make the `Create` function be `atomic` and thus lock this function.
	// 2. Count total submissions in system.
	// 3. Generate CPRN.
	// 4. Apply the CPRN to the submission.
	// 5. Unlock this `Create` function to be usable again by other calls.
	c.Logger.Debug("applying mutex")
	c.Kmutex.Lock("CPS-BACKEND-SUBMISSION-INSERTION") // Step 1
	defer func() {
		c.Kmutex.Unlock("CPS-BACKEND-SUBMISSION-INSERTION") // Step 5
		c.Logger.Debug("removing mutex")
	}()
	total, err := c.SubmissionStorer.CountAll(ctx) // Step 2
	if err != nil {
		c.Logger.Error("count all submissions error", slog.Any("error", err))
		return nil, err
	}
	m.CPSRN = c.CPSRN.GenerateNumber(userRole, total) // Step 3 & 4

	// DEVELOPERS NOTE:
	// Every submission creation is dependent on the `role` of the logged in
	// user in our system.
	switch userRole {
	case u_d.RetailerStaffRole:
		// Override state.
		m.State = s_d.SubmissionPendingState

		// Auto-assign the user-if
		m.UserFirstName = ctx.Value(constants.SessionUserFirstName).(string)
		m.UserLastName = ctx.Value(constants.SessionUserLastName).(string)
		m.ServiceType = s_d.PreScreeningServiceType
	case u_d.StaffRole:
		panic("SubmissionControllerImpl | Create | TODO: IMPLEMENT.")
		m.State = s_d.SubmissionActiveState
	default:
		m.State = s_d.SubmissionErrorState
	}

	// Update the `company name` field.
	userOrgID, ok := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	if ok {
		org, err := c.OrganizationStorer.GetByID(ctx, userOrgID)
		if err != nil {
			c.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if org == nil {
			c.Logger.Error("database get by id does not exist", slog.Any("organization id", userOrgID))
			return nil, fmt.Errorf("does not exist for organization id: %v", userOrgID)
		}
		m.OrganizationID = org.ID
		m.UserCompanyName = org.Name
	}

	// Add defaults.
	m.ID = primitive.NewObjectID()
	m.CreatedByUserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	m.CreatedAt = time.Now()
	m.ModifiedByUserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	m.ModifiedAt = time.Now()
	m.SubmissionDate = time.Now()

	// Attach a copy of the customer to our record.
	customerUser, err := c.UserStorer.GetByID(ctx, m.UserID)
	if err != nil {
		c.Logger.Error("database get customer by id error", slog.Any("error", err))
		return nil, err
	}
	if customerUser != nil {
		m.User = userToSubmissionUserCopy(customerUser)
	}

	// Save to our database.
	if err := c.SubmissionStorer.Create(ctx, m); err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	// Look up the publisher names and get the correct display name or get the other.
	var publisherNameDisplay string = constants.SubmissionPublisherNames[m.PublisherName]
	if m.PublisherName == constants.SubmissionPublisherNameOther {
		publisherNameDisplay = m.PublisherNameOther
	}

	// The next following lines of code will create the PDF file gnerator
	// request to be submitted into our PDF file generator to generate the data.
	r := &pdfbuilder.CBFFBuilderRequestDTO{
		CPSRN:                              m.CPSRN,
		Filename:                           fmt.Sprintf("%v.pdf", m.ID.Hex()),
		SubmissionDate:                     time.Now(),
		SeriesTitle:                        m.SeriesTitle,
		IssueVol:                           m.IssueVol,
		IssueNo:                            m.IssueNo,
		IssueCoverDate:                     m.IssueCoverDate,
		PublisherName:                      publisherNameDisplay,
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
	m.ModifiedAt = time.Now()

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
