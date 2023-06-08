package controller

import (
	"context"
	"errors"
	"fmt"
	go_os "os"
	"time"

	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	submission_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
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
	// Fetch the original submission.
	os, err := c.SubmissionStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, nil
	}

	// Modify our original submission.
	os.ModifiedAt = time.Now()
	// os.State = ns.State //BUGFIX - TODO WITH ROLES
	os.ServiceType = ns.ServiceType
	os.SubmissionDate = ns.SubmissionDate
	os.Item = ns.Item
	os.SeriesTitle = ns.SeriesTitle
	os.IssueVol = ns.IssueVol
	os.IssueNo = ns.IssueNo
	os.IssueCoverYear = ns.IssueCoverYear
	os.IssueCoverMonth = ns.IssueCoverMonth
	os.PublisherName = ns.PublisherName
	os.PublisherNameOther = ns.PublisherNameOther
	os.SpecialNotesLine1 = ns.SpecialNotesLine1
	os.SpecialNotesLine2 = ns.SpecialNotesLine2
	os.SpecialNotesLine3 = ns.SpecialNotesLine3
	os.SpecialNotesLine4 = ns.SpecialNotesLine4
	os.SpecialNotesLine5 = ns.SpecialNotesLine5
	os.GradingNotesLine1 = ns.GradingNotesLine1
	os.GradingNotesLine2 = ns.GradingNotesLine2
	os.GradingNotesLine3 = ns.GradingNotesLine3
	os.GradingNotesLine4 = ns.GradingNotesLine4
	os.GradingNotesLine5 = ns.GradingNotesLine5
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
	os.OverallNumberGrade = ns.OverallNumberGrade
	os.CpsPercentageGrade = ns.CpsPercentageGrade
	os.UserFirstName = ns.UserFirstName
	os.UserLastName = ns.UserLastName
	os.UserCompanyName = ns.UserCompanyName
	os.Filename = ns.Filename
	os.Item = fmt.Sprintf("%v, %v, %v", ns.SeriesTitle, ns.IssueVol, ns.IssueNo)

	// Save to the database the modified submission.
	if err := c.SubmissionStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	// Delete previous record from remote storage.
	if err := c.S3.DeleteByKeys(ctx, []string{os.FileUploadS3ObjectKey}); err != nil {
		c.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
		// Do not return an error, simply continue this function as there might
		// be a case were the file was removed on the s3 bucket by ourselves
		// or some other reason.
	}

	// Look up the publisher names and get the correct display name or get the other.
	var publisherNameDisplay string = constants.SubmissionPublisherNames[ns.PublisherName]
	if ns.PublisherName == constants.SubmissionPublisherNameOther {
		publisherNameDisplay = ns.PublisherNameOther
	}
	// The next following lines of code will create the PDF file gnerator
	// request to be submitted into our PDF file generator to generate the data.
	r := &pdfbuilder.CBFFBuilderRequestDTO{
		CPSRN:                              ns.CPSRN,
		Filename:                           fmt.Sprintf("%v.pdf", ns.CPSRN),
		SubmissionDate:                     time.Now(),
		SeriesTitle:                        ns.SeriesTitle,
		IssueVol:                           ns.IssueVol,
		IssueNo:                            ns.IssueNo,
		IssueCoverYear:                     ns.IssueCoverYear,
		IssueCoverMonth:                    ns.IssueCoverMonth,
		PublisherName:                      publisherNameDisplay,
		SpecialNotesLine1:                  ns.SpecialNotesLine1,
		SpecialNotesLine2:                  ns.SpecialNotesLine2,
		SpecialNotesLine3:                  ns.SpecialNotesLine3,
		SpecialNotesLine4:                  ns.SpecialNotesLine4,
		SpecialNotesLine5:                  ns.SpecialNotesLine5,
		GradingNotesLine1:                  ns.GradingNotesLine1,
		GradingNotesLine2:                  ns.GradingNotesLine2,
		GradingNotesLine3:                  ns.GradingNotesLine3,
		GradingNotesLine4:                  ns.GradingNotesLine4,
		GradingNotesLine5:                  ns.GradingNotesLine5,
		CreasesFinding:                     ns.CreasesFinding,
		TearsFinding:                       ns.TearsFinding,
		MissingPartsFinding:                ns.MissingPartsFinding,
		StainsFinding:                      ns.StainsFinding,
		DistortionFinding:                  ns.DistortionFinding,
		PaperQualityFinding:                ns.PaperQualityFinding,
		SpineFinding:                       ns.SpineFinding,
		CoverFinding:                       ns.CoverFinding,
		ShowsSignsOfTamperingOrRestoration: ns.ShowsSignsOfTamperingOrRestoration == 1,
		GradingScale:                       ns.GradingScale,
		OverallLetterGrade:                 ns.OverallLetterGrade,
		OverallNumberGrade:                 ns.OverallNumberGrade,
		CpsPercentageGrade:                 ns.CpsPercentageGrade,
		UserFirstName:                      ns.UserFirstName,
		UserLastName:                       ns.UserLastName,
		UserCompanyName:                    ns.UserCompanyName,
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
	if err := go_os.Remove(response.FilePath); err != nil {
		c.Logger.Warn("removing local file error", slog.Any("error", err))
		// Just continue even if we get an error...
	}

	return os, nil
}
