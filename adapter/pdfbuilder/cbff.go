package pdfbuilder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
)

type CBFFBuilderRequestDTO struct {
	ID                                 primitive.ObjectID `bson:"_id" json:"id"`
	Filename                           string             `bson:"filename" json:"filename"`
	SubmissionDate                     time.Time          `bson:"submission_date" json:"submission_date"`
	Item                               string             `bson:"item" json:"item"`
	SeriesTitle                        string             `bson:"series_title" json:"series_title"`
	IssueVol                           string             `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string             `bson:"issue_no" json:"issue_no"`
	IssueCoverDate                     string             `bson:"issue_cover_date" json:"issue_cover_date"`
	PublisherName                      string             `bson:"publisher_name" json:"publisher_name"`
	SpecialNotesLine1                  string             `bson:"special_notes_line_1" json:"special_notes_line_1"`
	SpecialNotesLine2                  string             `bson:"special_notes_line_2" json:"special_notes_line_2"`
	SpecialNotesLine3                  string             `bson:"special_notes_line_3" json:"special_notes_line_3"`
	SpecialNotesLine4                  string             `bson:"special_notes_line_4" json:"special_notes_line_4"`
	SpecialNotesLine5                  string             `bson:"special_notes_line_5" json:"special_notes_line_5"`
	GradingNotesLine1                  string             `bson:"grading_notes_line_1" json:"grading_notes_line_1"`
	GradingNotesLine2                  string             `bson:"grading_notes_line_2" json:"grading_notes_line_2"`
	GradingNotesLine3                  string             `bson:"grading_notes_line_3" json:"grading_notes_line_3"`
	GradingNotesLine4                  string             `bson:"grading_notes_line_4" json:"grading_notes_line_4"`
	GradingNotesLine5                  string             `bson:"grading_notes_line_5" json:"grading_notes_line_5"`
	CreasesFinding                     string             `bson:"creases_finding" json:"creases_finding"`
	TearsFinding                       string             `bson:"tears_finding" json:"tears_finding"`
	MissingPartsFinding                string             `bson:"missing_parts_finding" json:"missing_parts_finding"`
	StainsFinding                      string             `bson:"stains_finding" json:"stains_finding"`
	DistortionFinding                  string             `bson:"distortion_finding" json:"distortion_finding"`
	PaperQualityFinding                string             `bson:"paper_quality_finding" json:"paper_quality_finding"`
	SpineFinding                       string             `bson:"spine_finding" json:"spine_finding"`
	CoverFinding                       string             `bson:"cover_finding" json:"cover_finding"`
	ShowsSignsOfTamperingOrRestoration bool               `bson:"shows_signs_of_tampering_or_restoration" json:"shows_signs_of_tampering_or_restoration"`
	OverallLetterGrade                 string             `bson:"overall_letter_grade" json:"overall_letter_grade"`
	UserFirstName                      string             `bson:"user_first_name" json:"user_first_name"`
	UserLastName                       string             `bson:"user_last_name" json:"user_last_name"`
	UserCompanyName                    string             `bson:"user_company_name" json:"user_company_name"`
}
type CBFFBuilderResponseDTO struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
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

	pdf.SetFont("Helvetica", "B", 12)

	// ROW 1
	pdf.SetXY(113, 39)
	pdf.Cell(0, 0, fmt.Sprintf("%v", r.SubmissionDate.Day())) // Day
	pdf.SetXY(126, 39)
	pdf.Cell(0, 0, fmt.Sprintf("%v", int(r.SubmissionDate.Month()))) // Month (number)
	pdf.SetXY(135, 39)
	pdf.Cell(0, 0, fmt.Sprintf("%v", r.SubmissionDate.Year())) // Day

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
	// RIGHT
	//

	pdf.SetFont("Helvetica", "B", 14) // This controls the next text.

	// ROW 1 - Creases
	switch strings.ToLower(r.CreasesFinding) {
	case "pr":
		pdf.SetXY(92, 75)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(110, 75)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(127, 75)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(144, 75)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(163, 75)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(180, 75)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(197, 75)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, fmt.Errorf("missing value for crease finding with %v", r.CreasesFinding)
	}

	// ROW 2 - Tears
	switch strings.ToLower(r.TearsFinding) {
	case "pr":
		pdf.SetXY(92, 83)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(110, 83)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(127, 83)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(144, 83)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(163, 83)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(180, 83)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(197, 83)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for tears finding")
	}

	// ROW 3 - Missing Parts
	switch strings.ToLower(r.MissingPartsFinding) {
	case "pr":
		pdf.SetXY(92, 91)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(110, 91)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(127, 91)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(144, 91)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(163, 91)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(180, 91)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(197, 91)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for missing parts finding")
	}

	// ROW 4 - Stains / Marks / Substances
	switch strings.ToLower(r.StainsFinding) {
	case "pr":
		pdf.SetXY(92, 98)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(110, 98)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(127, 98)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(144, 98)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(163, 98)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(180, 98)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(197, 98)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for stains finding")
	}

	// ROW 5 - Distortion / Colour
	switch strings.ToLower(r.DistortionFinding) {
	case "pr":
		pdf.SetXY(92, 106)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(110, 106)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(127, 106)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(144, 106)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(163, 106)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(180, 106)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(197, 106)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for distorion finding")
	}

	// ROW 6 - Paper Quality
	switch strings.ToLower(r.PaperQualityFinding) {
	case "pr":
		pdf.SetXY(92, 113)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(110, 113)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(127, 113)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(144, 113)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(163, 113)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(180, 113)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(197, 113)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for paper quality finding")
	}

	// ROW 7 - Spine / Staples
	switch strings.ToLower(r.SpineFinding) {
	case "pr":
		pdf.SetXY(92, 121)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(110, 121)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(127, 121)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(144, 121)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(163, 121)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(180, 121)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(197, 121)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for paper quality finding")
	}

	// ROW 8 - Cover (Front & Back)
	switch strings.ToLower(r.CoverFinding) {
	case "pr":
		pdf.SetXY(92, 129)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(110, 129)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(127, 129)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(144, 129)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(163, 129)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(180, 129)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(197, 129)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value cover finding")
	}

	// ROW 9 - Shows signs of temp
	if r.ShowsSignsOfTamperingOrRestoration == true {
		pdf.SetXY(86, 136.5)
		pdf.Cell(0, 0, "X")
	} else {
		pdf.SetXY(101, 136.5)
		pdf.Cell(0, 0, "X")
	}

	pdf.SetFont("Helvetica", "B", 40)

	// ROW 10 - Grading
	pdf.SetXY(171, 153.5)
	pdf.Cell(0, 0, r.OverallLetterGrade)

	//
	// LEFT
	//

	pdf.SetFont("Helvetica", "B", 12)

	// ROW 1 - Special Notes
	pdf.SetXY(217, 75)
	pdf.Cell(0, 0, r.SpecialNotesLine1) // 17 characters.

	// ROW 2 - Special Notes
	pdf.SetXY(217, 83)
	pdf.Cell(0, 0, r.SpecialNotesLine2) // 17 characters.

	// ROW 3 - Special Notes
	pdf.SetXY(217, 91)
	pdf.Cell(0, 0, r.SpecialNotesLine3) // 17 characters.

	// ROW 4 - Special Notes
	pdf.SetXY(217, 99)
	pdf.Cell(0, 0, r.SpecialNotesLine4) // 17 characters.

	// ROW 5 - Special Notes
	pdf.SetXY(217, 107)
	pdf.Cell(0, 0, r.SpecialNotesLine5) // 17 characters.

	// ROW 1 - Grading Notes
	pdf.SetXY(217, 124)
	pdf.Cell(0, 0, r.GradingNotesLine1) // 17 characters.

	// ROW 2 - Grading Notes
	pdf.SetXY(217, 132)
	pdf.Cell(0, 0, r.GradingNotesLine2) // 17 characters.

	// ROW 3 - Grading Notes
	pdf.SetXY(217, 140)
	pdf.Cell(0, 0, r.GradingNotesLine3) // 17 characters.

	// ROW 4 - Grading Notes
	pdf.SetXY(217, 148)
	pdf.Cell(0, 0, r.GradingNotesLine4) // 17 characters.

	// ROW 5 - Grading Notes
	pdf.SetXY(217, 156)
	pdf.Cell(0, 0, r.GradingNotesLine5) // 17 characters.

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
		FilePath: filePath,
		Content:  bin,
	}, err
}
