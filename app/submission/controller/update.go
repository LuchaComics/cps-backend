package controller

import (
	"context"
	"errors"
	"fmt"
	go_os "os"
	"time"

	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	s_d "github.com/LuchaComics/cps-backend/app/submission/datastore"
	submission_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	u_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"
)

func (c *SubmissionControllerImpl) CreateComment(ctx context.Context, submissionID primitive.ObjectID, content string) (*submission_s.Submission, error) {
	// Fetch the original submission.
	s, err := c.SubmissionStorer.GetByID(ctx, submissionID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if s == nil {
		return nil, nil
	}

	// Create our comment.
	comment := &submission_s.SubmissionComment{
		ID:               primitive.NewObjectID(),
		Content:          content,
		OrganizationID:   ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID),
		CreatedByUserID:  ctx.Value(constants.SessionUserID).(primitive.ObjectID),
		CreatedByName:    ctx.Value(constants.SessionUserName).(string),
		CreatedAt:        time.Now(),
		ModifiedByUserID: ctx.Value(constants.SessionUserID).(primitive.ObjectID),
		ModifiedByName:   ctx.Value(constants.SessionUserName).(string),
		ModifiedAt:       time.Now(),
	}

	// Add our comment to the comments.
	s.ModifiedByUserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	s.ModifiedAt = time.Now()
	s.Comments = append(s.Comments, comment)

	// Save to the database the modified submission.
	if err := c.SubmissionStorer.UpdateByID(ctx, s); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return s, nil
}

func (c *SubmissionControllerImpl) SetUser(ctx context.Context, submissionID primitive.ObjectID, userID primitive.ObjectID) (*submission_s.Submission, error) {
	// Fetch the original submission.
	os, err := c.SubmissionStorer.GetByID(ctx, submissionID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, nil
	}

	// Fetch the original submission.
	cust, err := c.UserStorer.GetByID(ctx, userID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, nil
	}

	// Modify our original submission.
	os.ModifiedAt = time.Now()
	os.UserID = userID
	os.User = userToSubmissionUserCopy(cust)

	// Save to the database the modified submission.
	if err := c.SubmissionStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}

func (c *SubmissionControllerImpl) UpdateByID(ctx context.Context, ns *domain.Submission) (*domain.Submission, error) {
	// DEVELOPERS NOTE:
	// Every submission creation is dependent on the `role` of the logged in
	// user in our system so we need to extract it right away.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	//
	// Fetch submission.
	//

	// Fetch the original submission.
	os, err := c.SubmissionStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, nil
	}

	//
	// Set organization.
	//

	// DEVELOPERS NOTE:
	// Every submission creation is dependent on the `role` of the logged in
	// user in our system; however, the root administrator has the ability to
	// assign whatever organization you want.
	switch userRole {
	case u_d.UserRoleRetailer:
		c.Logger.Debug("retailer assigning their organization")
		os.OrganizationID = ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	case u_d.UserRoleRoot:
		c.Logger.Debug("admin picking custom organization")
	default:
		c.Logger.Error("unsupported role", slog.Any("role", userRole))
		return nil, fmt.Errorf("unsupported role via: %v", userRole)
	}

	// Lookup the organization.
	org, err := c.OrganizationStorer.GetByID(ctx, ns.OrganizationID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if org == nil {
		c.Logger.Error("database get by id does not exist", slog.Any("organization id", ns.OrganizationID))
		return nil, fmt.Errorf("does not exist for organization id: %v", ns.OrganizationID)
	}
	os.OrganizationID = org.ID
	os.OrganizationName = org.Name

	//
	// Update records in database.
	//

	// Modify our original submission.
	os.ModifiedAt = time.Now()
	// os.Status = ns.Status //BUGFIX - TODO WITH ROLES
	os.ServiceType = ns.ServiceType
	os.SubmissionDate = ns.SubmissionDate
	os.Item = fmt.Sprintf("%v, %v, %v", ns.SeriesTitle, ns.IssueVol, ns.IssueNo)
	os.SeriesTitle = ns.SeriesTitle
	os.IssueVol = ns.IssueVol
	os.IssueNo = ns.IssueNo
	os.IssueCoverYear = ns.IssueCoverYear
	os.IssueCoverMonth = ns.IssueCoverMonth
	os.PublisherName = ns.PublisherName
	os.PublisherNameOther = ns.PublisherNameOther
	os.SpecialNotes = ns.SpecialNotes
	os.GradingNotes = ns.GradingNotes
	os.CreasesFinding = ns.CreasesFinding
	os.TearsFinding = ns.TearsFinding
	os.MissingPartsFinding = ns.MissingPartsFinding
	os.StainsFinding = ns.StainsFinding
	os.DistortionFinding = ns.DistortionFinding
	os.PaperQualityFinding = ns.PaperQualityFinding
	os.SpineFinding = ns.SpineFinding
	os.CoverFinding = ns.CoverFinding
	os.ShowsSignsOfTamperingOrRestoration = ns.ShowsSignsOfTamperingOrRestoration
	os.GradingScale = ns.GradingScale
	os.OverallLetterGrade = ns.OverallLetterGrade
	os.IsOverallLetterGradeNearMintPlus = ns.IsOverallLetterGradeNearMintPlus
	os.OverallNumberGrade = ns.OverallNumberGrade
	os.CpsPercentageGrade = ns.CpsPercentageGrade
	// os.UserFirstName = ns.UserFirstName     // NO NEED TO CHANGE AFTER FACT.
	// os.UserLastName = ns.UserLastName       // NO NEED TO CHANGE AFTER FACT.
	// os.UserOrganizationName = ns.UserOrganizationName // NO NEED TO CHANGE AFTER FACT.
	os.Filename = ns.Filename
	os.Item = fmt.Sprintf("%v, %v, %v", ns.SeriesTitle, ns.IssueVol, ns.IssueNo)

	// Save to the database the modified submission.
	if err := c.SubmissionStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	//
	// Delete pdf file from s3
	//

	// Delete previous record from remote storage.
	if err := c.S3.DeleteByKeys(ctx, []string{os.FileUploadS3ObjectKey}); err != nil {
		c.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
		// Do not return an error, simply continue this function as there might
		// be a case were the file was removed on the s3 bucket by ourselves
		// or some other reason.
	}

	//
	// Create new PDF file.
	//

	// Look up the publisher names and get the correct display name or get the other.
	var publisherNameDisplay string = constants.SubmissionPublisherNames[ns.PublisherName]
	if ns.PublisherName == constants.SubmissionPublisherNameOther {
		publisherNameDisplay = ns.PublisherNameOther
	}

	pdfResponse := &pdfbuilder.PDFBuilderResponseDTO{}

	switch os.ServiceType {
	case s_d.ServiceTypePreScreening:
		// The next following lines of code will create the PDF file gnerator
		// request to be submitted into our PDF file generator to generate the data.
		r := &pdfbuilder.CBFFBuilderRequestDTO{
			CPSRN:                              os.CPSRN,
			Filename:                           fmt.Sprintf("%v.pdf", os.ID.Hex()),
			SubmissionDate:                     time.Now(),
			SeriesTitle:                        os.SeriesTitle,
			IssueVol:                           os.IssueVol,
			IssueNo:                            os.IssueNo,
			IssueCoverYear:                     os.IssueCoverYear,
			IssueCoverMonth:                    os.IssueCoverMonth,
			PublisherName:                      publisherNameDisplay,
			SpecialNotes:                       os.SpecialNotes,
			GradingNotes:                       os.GradingNotes,
			CreasesFinding:                     os.CreasesFinding,
			TearsFinding:                       os.TearsFinding,
			MissingPartsFinding:                os.MissingPartsFinding,
			StainsFinding:                      os.StainsFinding,
			DistortionFinding:                  os.DistortionFinding,
			PaperQualityFinding:                os.PaperQualityFinding,
			SpineFinding:                       os.SpineFinding,
			CoverFinding:                       os.CoverFinding,
			ShowsSignsOfTamperingOrRestoration: os.ShowsSignsOfTamperingOrRestoration == 1,
			GradingScale:                       os.GradingScale,
			OverallLetterGrade:                 os.OverallLetterGrade,
			IsOverallLetterGradeNearMintPlus:   os.IsOverallLetterGradeNearMintPlus,
			OverallNumberGrade:                 os.OverallNumberGrade,
			CpsPercentageGrade:                 os.CpsPercentageGrade,
			UserFirstName:                      os.UserFirstName,
			UserLastName:                       os.UserLastName,
			UserOrganizationName:               os.OrganizationName,
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
			CPSRN:                              os.CPSRN,
			Filename:                           fmt.Sprintf("%v.pdf", os.ID.Hex()),
			SubmissionDate:                     time.Now(),
			SeriesTitle:                        os.SeriesTitle,
			IssueVol:                           os.IssueVol,
			IssueNo:                            os.IssueNo,
			IssueCoverYear:                     os.IssueCoverYear,
			IssueCoverMonth:                    os.IssueCoverMonth,
			PublisherName:                      publisherNameDisplay,
			SpecialNotes:                       os.SpecialNotes,
			GradingNotes:                       os.GradingNotes,
			CreasesFinding:                     os.CreasesFinding,
			TearsFinding:                       os.TearsFinding,
			MissingPartsFinding:                os.MissingPartsFinding,
			StainsFinding:                      os.StainsFinding,
			DistortionFinding:                  os.DistortionFinding,
			PaperQualityFinding:                os.PaperQualityFinding,
			SpineFinding:                       os.SpineFinding,
			CoverFinding:                       os.CoverFinding,
			ShowsSignsOfTamperingOrRestoration: os.ShowsSignsOfTamperingOrRestoration == 1,
			GradingScale:                       os.GradingScale,
			OverallLetterGrade:                 os.OverallLetterGrade,
			OverallNumberGrade:                 os.OverallNumberGrade,
			CpsPercentageGrade:                 os.CpsPercentageGrade,
			UserFirstName:                      os.UserFirstName,
			UserLastName:                       os.UserLastName,
			UserOrganizationName:               os.OrganizationName,
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

	//
	// Update record with PDF file information with record.
	//

	// The following will save the S3 key of our file upload into our record.
	os.FileUploadS3ObjectKey = path
	os.ModifiedAt = time.Now()

	if err := c.SubmissionStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update error", slog.Any("error", err))
		return nil, err
	}

	// The following will generate a pre-signed URL so user can download the file.
	downloadableURL, err := c.S3.GetDownloadablePresignedURL(ctx, os.FileUploadS3ObjectKey, time.Minute*15)
	if err != nil {
		c.Logger.Error("s3 presign error", slog.Any("error", err))
		return nil, err
	}
	os.FileUploadDownloadableFileURL = downloadableURL

	// Removing local file from the directory and don't do anything if we have errors.
	if err := go_os.Remove(pdfResponse.FilePath); err != nil {
		c.Logger.Warn("removing local file error", slog.Any("error", err))
		// Just continue even if we get an error...
	}

	return os, nil
}
