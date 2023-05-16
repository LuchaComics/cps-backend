package pdfbuilder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
)

type CBFFBuilderRequestDTO struct {
	ID                       primitive.ObjectID `bson:"_id" json:"id"`
	Date                     time.Time          `bson:"date" json:"date"`
	Item                     string             `bson:"item" json:"item"` // Created by system.
	SeriesTitle              string             `bson:"series_title" json:"series_title"`
	IssueVol                 string             `bson:"issue_vol" json:"issue_vol"`
	IssueNo                  string             `bson:"issue_no" json:"issue_no"`
	IssueCoverDate           string             `bson:"issue_cover_date" json:"issue_cover_date"`
	PublisherName            string             `bson:"publisher_name" json:"publisher_name"`
	IssueSpecialDetails      string             `bson:"issue_special_details" json:"issue_special_details"`
	CreasesFinding           string             `bson:"creases_finding" json:"creases_finding"`
	TearsFinding             string             `bson:"tears_finding" json:"tears_finding"`
	MissingPartsFinding      string             `bson:"missing_parts_finding" json:"missing_parts_finding"`
	StainsFinding            string             `bson:"stains_finding" json:"stains_finding"`
	DistortionFinding        string             `bson:"distortion_finding" json:"distortion_finding"`
	PaperQualityFinding      string             `bson:"paper_quality_finding" json:"paper_quality_finding"`
	SpineFinding             string             `bson:"spine_finding" json:"spine_finding"`
	CoverFinding             string             `bson:"cover_finding" json:"cover_finding"`
	OtherFinding             string             `bson:"other_finding" json:"other_finding"`
	OtherFindingText         string             `bson:"other_finding_text" json:"other_finding_text"`
	OverallLetterGrade       string             `bson:"overall_letter_grade" json:"overall_letter_grade"`
	UserID                   primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserFirstName            string             `bson:"user_first_name" json:"user_first_name"`
	UserLastName             string             `bson:"user_last_name" json:"user_last_name"`
	UserCompanyName          string             `bson:"user_company_name" json:"user_company_name"`
	UserSignature            string             `bson:"user_signature" json:"user_signature"`
	InspectorSignature       string             `bson:"inspector_signature" json:"inspector_signature"`
	InspectorDate            time.Time          `bson:"inspector_date" json:"inspector_date"`
	InspectorFirstName       string             `bson:"inspector_first_name" json:"inspector_first_name"`
	InspectorLastName        string             `bson:"inspector_last_name" json:"inspector_last_name"`
	InspectorCompany         string             `bson:"inspector_company_name" json:"inspector_company_name"`
	SecondInspectorSignature string             `bson:"second_inspector_signature" json:"second_inspector_signature"`
	SecondInspectorFirstName string             `bson:"second_inspector_first_name" json:"second_inspector_first_name"`
	SecondInspectorLastName  string             `bson:"second_inspector_last_name" json:"second_inspector_last_name"`
	SecondInspectorCompany   string             `bson:"second_inspector_company" json:"second_inspector_company"`
	SecondInspectorDate      time.Time          `bson:"second_inspector_date" json:"second_inspector_date"`
	ThirdInspectorSignature  string             `bson:"third_inspector_signature" json:"third_inspector_signature"`
	ThirdInspectorFirstName  string             `bson:"third_inspector_first_name" json:"third_inspector_first_name"`
	ThirdInspectorLastName   string             `bson:"third_inspector_last_name" json:"third_inspector_last_name"`
	ThirdInspectorCompany    string             `bson:"third_inspector_company" json:"third_inspector_company"`
	ThirdInspectorDate       time.Time          `bson:"third_inspector_date" json:"third_inspector_date"`
}
type CBFFBuilderResponseDTO struct {
	FileName string `json:"file_name"`
	Content  []byte `json:"content"`
}

type CBFFBuilder interface {
	GeneratePDF(dto *CBFFBuilderRequestDTO) (*CBFFBuilderResponseDTO, error)
}

type cbffBuilder struct {
	PDFTemplateFilePath string
	DataDirectoryPath   string
	UUID                uuid.Provider
	Logger              *slog.Logger
}

func NewCBFFBuilder(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) CBFFBuilder {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("pdf builder for cbff initializing...")
	_, err := os.Stat(cfg.PDFBuilder.CBFFTemplatePath)
	if os.IsNotExist(err) {
		log.Fatal(errors.New("file does not exist"))
	}

	return &cbffBuilder{
		PDFTemplateFilePath: cfg.PDFBuilder.CBFFTemplatePath,
		DataDirectoryPath:   cfg.PDFBuilder.DataDirectoryPath,
		UUID:                uuidp,
		Logger:              logger,
	}
}

func (bdr *cbffBuilder) GeneratePDF(r *CBFFBuilderRequestDTO) (*CBFFBuilderResponseDTO, error) {
	var err error

	// Open our PDF invoice template and create clone it for the PDF invoice we will be building with.
	pdf := gofpdf.New("L", "mm", "A4", "")
	tpl1 := gofpdi.ImportPage(pdf, bdr.PDFTemplateFilePath, 1, "/MediaBox")

	pdf.AddPage()

	// Draw imported template onto page
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 0, 297, 210)

	//
	// UNIQUE IDENTIFIER
	//

	pdf.SetFont("Courier", "", 12)
	pdf.SetXY(17, 21)
	pdf.Cell(0, 0, r.ID.Hex())

	//
	// LEFT SIDE
	//

	pdf.SetFont("Helvetica", "", 12)

	// ROW 1
	pdf.SetXY(113, 39)
	pdf.Cell(0, 0, fmt.Sprintf("%v", r.Date.Day())) // Day
	pdf.SetXY(126, 39)
	pdf.Cell(0, 0, fmt.Sprintf("%v", int(r.Date.Month()))) // Month (number)
	pdf.SetXY(135, 39)
	pdf.Cell(0, 0, fmt.Sprintf("%v", r.Date.Year())) // Day

	// ROW 2
	pdf.SetXY(82, 47)
	pdf.Cell(0, 0, r.UserFirstName)
	pdf.SetXY(114, 47)
	pdf.Cell(0, 0, r.UserLastName)

	// ROW 3
	pdf.SetXY(27, 56)
	pdf.Cell(0, 0, r.UserCompanyName)

	//
	// RIGHT SIDE
	//

	// ROW 1
	pdf.SetXY(162, 39)
	pdf.Cell(0, 0, r.SeriesTitle)

	// ROW 2
	pdf.SetXY(160, 47.5)
	pdf.Cell(0, 0, r.IssueVol)
	pdf.SetXY(193, 47.5)
	pdf.Cell(0, 0, r.IssueNo)
	pdf.SetXY(235, 47.5)
	pdf.Cell(0, 0, r.IssueCoverDate)

	// ROW 3
	pdf.SetXY(172, 56)
	pdf.Cell(0, 0, r.PublisherName)

	//
	// MIDDLE
	//

	pdf.SetFont("Helvetica", "B", 14) // This controls the next text.

	// ROW 1
	pdf.SetXY(92, 75)
	pdf.Cell(0, 0, "PR")
	pdf.SetXY(110, 75)
	pdf.Cell(0, 0, "FR")
	pdf.SetXY(127, 75)
	pdf.Cell(0, 0, "GD")
	pdf.SetXY(144, 75)
	pdf.Cell(0, 0, "VG")
	pdf.SetXY(163, 75)
	pdf.Cell(0, 0, "FN")
	pdf.SetXY(180, 75)
	pdf.Cell(0, 0, "VF")
	pdf.SetXY(197, 75)
	pdf.Cell(0, 0, "NM")

	// ROW 2
	pdf.SetXY(92, 83)
	pdf.Cell(0, 0, "PR")
	pdf.SetXY(110, 83)
	pdf.Cell(0, 0, "FR")
	pdf.SetXY(127, 83)
	pdf.Cell(0, 0, "GD")
	pdf.SetXY(144, 83)
	pdf.Cell(0, 0, "VG")
	pdf.SetXY(163, 83)
	pdf.Cell(0, 0, "FN")
	pdf.SetXY(180, 83)
	pdf.Cell(0, 0, "VF")
	pdf.SetXY(197, 83)
	pdf.Cell(0, 0, "NM")

	// ROW 3
	pdf.SetXY(92, 91)
	pdf.Cell(0, 0, "PR")
	pdf.SetXY(110, 91)
	pdf.Cell(0, 0, "FR")
	pdf.SetXY(127, 91)
	pdf.Cell(0, 0, "GD")
	pdf.SetXY(144, 91)
	pdf.Cell(0, 0, "VG")
	pdf.SetXY(163, 91)
	pdf.Cell(0, 0, "FN")
	pdf.SetXY(180, 91)
	pdf.Cell(0, 0, "VF")
	pdf.SetXY(197, 91)
	pdf.Cell(0, 0, "NM")

	// ROW 4
	pdf.SetXY(92, 98)
	pdf.Cell(0, 0, "PR")
	pdf.SetXY(110, 98)
	pdf.Cell(0, 0, "FR")
	pdf.SetXY(127, 98)
	pdf.Cell(0, 0, "GD")
	pdf.SetXY(144, 98)
	pdf.Cell(0, 0, "VG")
	pdf.SetXY(163, 98)
	pdf.Cell(0, 0, "FN")
	pdf.SetXY(180, 98)
	pdf.Cell(0, 0, "VF")
	pdf.SetXY(197, 98)
	pdf.Cell(0, 0, "NM")

	// ROW 5
	pdf.SetXY(92, 106)
	pdf.Cell(0, 0, "PR")
	pdf.SetXY(110, 106)
	pdf.Cell(0, 0, "FR")
	pdf.SetXY(127, 106)
	pdf.Cell(0, 0, "GD")
	pdf.SetXY(144, 106)
	pdf.Cell(0, 0, "VG")
	pdf.SetXY(163, 106)
	pdf.Cell(0, 0, "FN")
	pdf.SetXY(180, 106)
	pdf.Cell(0, 0, "VF")
	pdf.SetXY(197, 106)
	pdf.Cell(0, 0, "NM")

	// ROW 6
	pdf.SetXY(92, 113)
	pdf.Cell(0, 0, "PR")
	pdf.SetXY(110, 113)
	pdf.Cell(0, 0, "FR")
	pdf.SetXY(127, 113)
	pdf.Cell(0, 0, "GD")
	pdf.SetXY(144, 113)
	pdf.Cell(0, 0, "VG")
	pdf.SetXY(163, 113)
	pdf.Cell(0, 0, "FN")
	pdf.SetXY(180, 113)
	pdf.Cell(0, 0, "VF")
	pdf.SetXY(197, 113)
	pdf.Cell(0, 0, "NM")

	// ROW 7
	pdf.SetXY(92, 121)
	pdf.Cell(0, 0, "PR")
	pdf.SetXY(110, 121)
	pdf.Cell(0, 0, "FR")
	pdf.SetXY(127, 121)
	pdf.Cell(0, 0, "GD")
	pdf.SetXY(144, 121)
	pdf.Cell(0, 0, "VG")
	pdf.SetXY(163, 121)
	pdf.Cell(0, 0, "FN")
	pdf.SetXY(180, 121)
	pdf.Cell(0, 0, "VF")
	pdf.SetXY(197, 121)
	pdf.Cell(0, 0, "NM")

	// ROW 8
	pdf.SetXY(92, 129)
	pdf.Cell(0, 0, "PR")
	pdf.SetXY(110, 129)
	pdf.Cell(0, 0, "FR")
	pdf.SetXY(127, 129)
	pdf.Cell(0, 0, "GD")
	pdf.SetXY(144, 129)
	pdf.Cell(0, 0, "VG")
	pdf.SetXY(163, 129)
	pdf.Cell(0, 0, "FN")
	pdf.SetXY(180, 129)
	pdf.Cell(0, 0, "VF")
	pdf.SetXY(197, 129)
	pdf.Cell(0, 0, "NM")

	////
	//// Generate the file and save it to the file.
	////

	fileName := fmt.Sprintf("%s.pdf", bdr.UUID.NewUUID())
	filePath := fmt.Sprintf("%s/%s", bdr.DataDirectoryPath, fileName)

	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		return nil, err
	}

	////
	//// Open the file and read all the binary data.
	////

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bin, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	////
	//// Return the generate invoice.
	////

	return &CBFFBuilderResponseDTO{
		FileName: fileName,
		Content:  bin,
	}, err
}
