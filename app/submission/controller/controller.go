package controller

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	domain "github.com/LuchaComics/cps-backend/app/submission/datastore"
	submission_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// SubmissionController Interface for submission business logic controller.
type SubmissionController interface {
	Create(ctx context.Context, m *domain.Submission) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Submission, error)
	UpdateByID(ctx context.Context, m *domain.Submission) error
	ListByFilter(ctx context.Context, f *domain.SubmissionListFilter) (*domain.SubmissionListResult, error)
}

type SubmissionControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	Password         password.Provider
	CBFFBuilder      pdfbuilder.CBFFBuilder
	SubmissionStorer submission_s.SubmissionStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	cbffb pdfbuilder.CBFFBuilder,
	sub_storer submission_s.SubmissionStorer,
) SubmissionController {

	// FOR TESTING PURPOSES ONLY.
	r := &pdfbuilder.CBFFBuilderRequestDTO{
		ID:   primitive.NewObjectID(),
		Date: time.Now(),
		// CreatedTime              time.Time          `bson:"created_time,omitempty" json:"created_time,omitempty"`
		// ModifiedTime             time.Time          `bson:"modified_time,omitempty" json:"modified_time,omitempty"`
		// ServiceType              int8               `bson:"service_type" json:"service_type"`
		// State                    int8               `bson:"state" json:"state"`
		// Item                     string             `bson:"item" json:"item"` // Created by system.
		SeriesTitle:    "Winter World",
		IssueVol:       "Vol 1",
		IssueNo:        "#1",
		IssueCoverDate: "16/05/2000",
		PublisherName:  "Some publisher",
		// IssueSpecialDetails      string             `bson:"issue_special_details" json:"issue_special_details"`
		// CreasesFinding           string             `bson:"creases_finding" json:"creases_finding"`
		// TearsFinding             string             `bson:"tears_finding" json:"tears_finding"`
		// MissingPartsFinding      string             `bson:"missing_parts_finding" json:"missing_parts_finding"`
		// StainsFinding            string             `bson:"stains_finding" json:"stains_finding"`
		// DistortionFinding        string             `bson:"distortion_finding" json:"distortion_finding"`
		// PaperQualityFinding      string             `bson:"paper_quality_finding" json:"paper_quality_finding"`
		// SpineFinding             string             `bson:"spine_finding" json:"spine_finding"`
		// CoverFinding             string             `bson:"cover_finding" json:"cover_finding"`
		// OtherFinding             string             `bson:"other_finding" json:"other_finding"`
		// OtherFindingText         string             `bson:"other_finding_text" json:"other_finding_text"`
		// OverallLetterGrade       string             `bson:"overall_letter_grade" json:"overall_letter_grade"`
		// UserID                   primitive.ObjectID `bson:"user_id" json:"user_id"`
		UserFirstName:   "Bartlomiej",
		UserLastName:    "Miks",
		UserCompanyName: "Mika Software Corporation",
		// UserSignature            string             `bson:"user_signature" json:"user_signature"`
		// InspectorSignature       string             `bson:"inspector_signature" json:"inspector_signature"`
		// InspectorDate            time.Time          `bson:"inspector_date" json:"inspector_date"`
		// InspectorFirstName       string             `bson:"inspector_first_name" json:"inspector_first_name"`
		// InspectorLastName        string             `bson:"inspector_last_name" json:"inspector_last_name"`
		// InspectorCompany         string             `bson:"inspector_company_name" json:"inspector_company_name"`
		// SecondInspectorSignature string             `bson:"second_inspector_signature" json:"second_inspector_signature"`
		// SecondInspectorFirstName string             `bson:"second_inspector_first_name" json:"second_inspector_first_name"`
		// SecondInspectorLastName  string             `bson:"second_inspector_last_name" json:"second_inspector_last_name"`
		// SecondInspectorCompany   string             `bson:"second_inspector_company" json:"second_inspector_company"`
		// SecondInspectorDate      time.Time          `bson:"second_inspector_date" json:"second_inspector_date"`
		// ThirdInspectorSignature  string             `bson:"third_inspector_signature" json:"third_inspector_signature"`
		// ThirdInspectorFirstName  string             `bson:"third_inspector_first_name" json:"third_inspector_first_name"`
		// ThirdInspectorLastName   string             `bson:"third_inspector_last_name" json:"third_inspector_last_name"`
		// ThirdInspectorCompany    string             `bson:"third_inspector_company" json:"third_inspector_company"`
		// ThirdInspectorDate       time.Time          `bson:"third_inspector_date" json:"third_inspector_date"`
	}
	res, err := cbffb.GeneratePDF(r)
	log.Println("===--->", res, err, "<---===")

	s := &SubmissionControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		Password:         passwordp,
		CBFFBuilder:      cbffb,
		SubmissionStorer: sub_storer,
	}
	s.Logger.Debug("submission controller initialization started...")
	s.Logger.Debug("submission controller initialized")
	return s
}
