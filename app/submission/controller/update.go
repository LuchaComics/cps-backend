package controller

import (
	"context"
	"time"

	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"golang.org/x/exp/slog"
)

func (c *SubmissionControllerImpl) UpdateByID(ctx context.Context, ns *domain.Submission) error {
	// Fetch the original submission.
	os, err := c.SubmissionStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if os == nil {
		return nil
	}

	// Modify our original submission.
	os.ModifiedTime = time.Now()
	os.ServiceType = ns.ServiceType
	os.State = ns.State
	os.Item = ns.Item
	os.SeriesTitle = ns.SeriesTitle
	os.IssueVol = ns.IssueVol
	os.IssueNo = ns.IssueNo
	os.IssueCoverDate = ns.IssueCoverDate
	os.IssueSpecialDetails = ns.IssueSpecialDetails
	os.CreasesFinding = ns.CreasesFinding
	os.TearsFinding = ns.TearsFinding
	os.MissingPartsFinding = ns.MissingPartsFinding
	os.StainsFinding = ns.StainsFinding
	os.DistortionFinding = ns.DistortionFinding
	os.PaperQualityFinding = ns.PaperQualityFinding
	os.SpineFinding = ns.SpineFinding
	os.CoverFinding = ns.CoverFinding
	os.OtherFinding = ns.OtherFinding
	os.OtherFindingText = ns.OtherFindingText
	os.OverallLetterGrade = ns.OverallLetterGrade

	// Save to the database the modified submission.
	if err := c.SubmissionStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return err
}
