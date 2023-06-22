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
	domain "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	s_d "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	u_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
)

func (c *ComicSubmissionControllerImpl) Create(ctx context.Context, m *s_d.ComicSubmission) (*s_d.ComicSubmission, error) {
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
	f := &domain.ComicSubmissionListFilter{CreatedByUserRole: userRole} // Part of ID requires count of staff or retailer.
	total, err := c.ComicSubmissionStorer.CountByFilter(ctx, f)         // Step 2
	if err != nil {
		c.Logger.Error("count all submissions error", slog.Any("error", err))
		return nil, err
	}
	m.CPSRN = c.CPSRN.GenerateNumber(userRole, total) // Step 3 & 4
	c.Logger.Error("Generated CPSRN",
		slog.String("CPSRN", m.CPSRN),
		slog.Int64("Role", int64(userRole)),
		slog.Int64("total", total))

	// Auto-assign the user-if
	m.UserFirstName = ctx.Value(constants.SessionUserFirstName).(string)
	m.UserLastName = ctx.Value(constants.SessionUserLastName).(string)

	// DEVELOPERS NOTE:
	// Every submission creation is dependent on the `role` of the logged in
	// user in our system; however, the root administrator has the ability to
	// assign whatever organization you want.
	switch userRole {
	case u_d.UserRoleRetailer:
		c.Logger.Debug("retailer assigning their organization")
		m.OrganizationID = ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	case u_d.UserRoleRoot:
		c.Logger.Debug("admin picking custom organization")
	default:
		c.Logger.Error("unsupported role", slog.Any("role", userRole))
		return nil, fmt.Errorf("unsupported role via: %v", userRole)
	}

	// Lookup the organization.
	org, err := c.OrganizationStorer.GetByID(ctx, m.OrganizationID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if org == nil {
		c.Logger.Error("database get by id does not exist", slog.Any("organization id", m.OrganizationID))
		return nil, fmt.Errorf("does not exist for organization id: %v", m.OrganizationID)
	}
	m.OrganizationID = org.ID
	m.OrganizationName = org.Name

	// Add defaults.
	m.ID = primitive.NewObjectID()
	m.CreatedByUserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	m.CreatedByUserRole = userRole
	m.CreatedAt = time.Now()
	m.ModifiedByUserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	m.ModifiedByUserRole = userRole
	m.ModifiedAt = time.Now()
	m.SubmissionDate = time.Now()
	m.Item = fmt.Sprintf("%v, %v, %v", m.SeriesTitle, m.IssueVol, m.IssueNo)

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
	if err := c.ComicSubmissionStorer.Create(ctx, m); err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	// Look up the publisher names and get the correct display name or get the other.
	var publisherNameDisplay string = constants.SubmissionPublisherNames[m.PublisherName]
	if m.PublisherName == constants.SubmissionPublisherNameOther {
		publisherNameDisplay = m.PublisherNameOther
	}

	//
	// Generate the PDF file based on the `service type`.
	//

	pdfResponse := &pdfbuilder.PDFBuilderResponseDTO{}

	switch m.ServiceType {
	case s_d.ServiceTypePreScreening:
		// The next following lines of code will create the PDF file gnerator
		// request to be submitted into our PDF file generator to generate the data.
		r := &pdfbuilder.CBFFBuilderRequestDTO{
			CPSRN:                              m.CPSRN,
			Filename:                           fmt.Sprintf("%v.pdf", m.ID.Hex()),
			SubmissionDate:                     time.Now(),
			SeriesTitle:                        m.SeriesTitle,
			IssueVol:                           m.IssueVol,
			IssueNo:                            m.IssueNo,
			IssueCoverYear:                     m.IssueCoverYear,
			IssueCoverMonth:                    m.IssueCoverMonth,
			PublisherName:                      publisherNameDisplay,
			SpecialNotes:                       m.SpecialNotes,
			GradingNotes:                       m.GradingNotes,
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
			UserOrganizationName:               m.OrganizationName,
		}
		pdfResponse, err = c.CBFFBuilder.GeneratePDF(r)
		if err != nil {
			c.Logger.Error("generate pdf error", slog.Any("error", err))
			return nil, err
		}
		if pdfResponse == nil {
			c.Logger.Error("generate pdf error does not return a response")
			return nil, errors.New("no response from pdf generator")
		}
	case s_d.ServiceTypePedigree:
		// The next following lines of code will create the PDF file gnerator
		// request to be submitted into our PDF file generator to generate the data.
		r := &pdfbuilder.PCBuilderRequestDTO{
			CPSRN:                              m.CPSRN,
			Filename:                           fmt.Sprintf("%v.pdf", m.ID.Hex()),
			SubmissionDate:                     time.Now(),
			SeriesTitle:                        m.SeriesTitle,
			IssueVol:                           m.IssueVol,
			IssueNo:                            m.IssueNo,
			IssueCoverYear:                     m.IssueCoverYear,
			IssueCoverMonth:                    m.IssueCoverMonth,
			PublisherName:                      publisherNameDisplay,
			SpecialNotes:                       m.SpecialNotes,
			GradingNotes:                       m.GradingNotes,
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
			UserOrganizationName:               m.OrganizationName,
		}
		pdfResponse, err = c.PCBuilder.GeneratePDF(r)
		if err != nil {
			c.Logger.Error("generate pdf error", slog.Any("error", err))
			return nil, err
		}
		if pdfResponse == nil {
			c.Logger.Error("generate pdf error does not return a response")
			return nil, errors.New("no response from pdf generator")
		}
	}

	// The next few lines will upload our PDF to our remote storage. Once the
	// file is saved remotely, we will have a connection to it through a "key"
	// unique reference to the uploaded file.
	path := fmt.Sprintf("uploads/%v", pdfResponse.FileName)
	err = c.S3.UploadContent(ctx, path, pdfResponse.Content)
	if err != nil {
		c.Logger.Error("s3 upload error", slog.Any("error", err))
		return nil, err
	}

	// The following will save the S3 key of our file upload into our record.
	m.FileUploadS3ObjectKey = path
	m.ModifiedAt = time.Now()

	if err := c.ComicSubmissionStorer.UpdateByID(ctx, m); err != nil {
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
	if err := os.Remove(pdfResponse.FilePath); err != nil {
		c.Logger.Warn("removing local file error", slog.Any("error", err))
		// Just continue even if we get an error...
	}

	// The following code will send the email notifications to the correct individuals.
	if err := c.sendNewComicSubmissionEmails(m); err != nil {
		c.Logger.Error("database update error", slog.Any("error", err))
		// Do not return error, just keep it in the server logs.
	}
	return m, nil
}