package pdfbuilder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// "strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
	"golang.org/x/exp/slog"

	s_d "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	c "github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

type CCBuilderRequestDTO struct {
	CPSRN                              string                     `bson:"cpsrn" json:"cpSrn"`
	Filename                           string                     `bson:"filename" json:"filename"`
	SubmissionDate                     time.Time                  `bson:"submission_date" json:"submission_date"`
	Item                               string                     `bson:"item" json:"item"`
	SeriesTitle                        string                     `bson:"series_title" json:"series_title"`
	IssueVol                           string                     `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string                     `bson:"issue_no" json:"issue_no"`
	IssueCoverYear                     int64                      `bson:"issue_cover_year" json:"issue_cover_year"`
	IssueCoverMonth                    int8                       `bson:"issue_cover_month" json:"issue_cover_month"`
	PublisherName                      string                     `bson:"publisher_name" json:"publisher_name"`
	SpecialNotes                       string                     `bson:"special_notes" json:"special_notes"`
	GradingNotes                       string                     `bson:"grading_notes" json:"grading_notes"`
	CreasesFinding                     string                     `bson:"creases_finding" json:"creases_finding"`
	TearsFinding                       string                     `bson:"tears_finding" json:"tears_finding"`
	MissingPartsFinding                string                     `bson:"missing_parts_finding" json:"missing_parts_finding"`
	StainsFinding                      string                     `bson:"stains_finding" json:"stains_finding"`
	DistortionFinding                  string                     `bson:"distortion_finding" json:"distortion_finding"`
	PaperQualityFinding                string                     `bson:"paper_quality_finding" json:"paper_quality_finding"`
	SpineFinding                       string                     `bson:"spine_finding" json:"spine_finding"`
	CoverFinding                       string                     `bson:"cover_finding" json:"cover_finding"`
	ShowsSignsOfTamperingOrRestoration bool                       `bson:"shows_signs_of_tampering_or_restoration" json:"shows_signs_of_tampering_or_restoration"`
	GradingScale                       int8                       `bson:"grading_scale" json:"grading_scale"`
	OverallLetterGrade                 string                     `bson:"overall_letter_grade" json:"overall_letter_grade"`
	IsOverallLetterGradeNearMintPlus   bool                       `bson:"is_overall_letter_grade_near_mint_plus" json:"is_overall_letter_grade_near_mint_plus"`
	IsCpsIndieMintGem                  bool                       `bson:"is_cps_indie_mint_gem" json:"is_cps_indie_mint_gem"`
	OverallNumberGrade                 float64                    `bson:"overall_number_grade" json:"overall_number_grade"`
	CpsPercentageGrade                 float64                    `bson:"cps_percentage_grade" json:"cps_percentage_grade"`
	UserFirstName                      string                     `bson:"user_first_name" json:"user_first_name"`
	UserLastName                       string                     `bson:"user_last_name" json:"user_last_name"`
	UserOrganizationName               string                     `bson:"user_organization_name" json:"user_organization_name"`
	Signatures                         []*s_d.SubmissionSignature `bson:"signatures" json:"signatures,omitempty"`
	SpecialDetails                     int8                       `bson:"special_details" json:"special_details"`
	SpecialDetailsOther                string                     `bson:"special_details_other" json:"special_details_other"`
}

// CCBuilder interface for building the "CPS C-Capsule Indie Mint Gem" edition document.
type CCBuilder interface {
	GeneratePDF(dto *CCBuilderRequestDTO) (*PDFBuilderResponseDTO, error)
}

type ccBuilder struct {
	PDFTemplateFilePath string
	DataDirectoryPath   string
	UUID                uuid.Provider
	Logger              *slog.Logger
}

func NewCCBuilder(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) CCBuilder {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("pdf builder for CC initializing...")
	_, err := os.Stat(cfg.PDFBuilder.CCTemplatePath)
	if os.IsNotExist(err) {
		log.Fatal(errors.New("file does not exist"))
	}

	return &ccBuilder{
		PDFTemplateFilePath: cfg.PDFBuilder.CCTemplatePath,
		DataDirectoryPath:   cfg.PDFBuilder.DataDirectoryPath,
		UUID:                uuidp,
		Logger:              logger,
	}
}

func (bdr *ccBuilder) GeneratePDF(r *CCBuilderRequestDTO) (*PDFBuilderResponseDTO, error) {
	var err error

	// Open our PDF invoice template and create clone it for the PDF invoice we will be building with.
	pdf := gofpdf.New("P", "mm", "A4", "")
	tpl1 := gofpdi.ImportPage(pdf, bdr.PDFTemplateFilePath, 1, "/MediaBox")

	pdf.AddPage()

	// Draw imported template onto page
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 0, 210, 300)

	//
	// CPS REGISTRY NUMBER
	//

	pdf.SetFont("Courier", "", 12)

	// Set the transformation matrix to rotate 180 degrees
	pdf.TransformBegin()
	pdf.TransformRotate(180, 190, 27) // angle=180, x=190, y=27

	// Print the text
	pdf.Text(190, 27, r.CPSRN) // x=190, y=27

	pdf.TransformEnd()

	//
	// TITLE
	//

	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetXY(80, 51)
	pdf.Cell(0, 0, r.PublisherName)

	//
	// LEFT SIDE
	//

	pdf.SetFont("Helvetica", "B", 14)

	// ROW 1
	pdf.SetXY(60, 60)
	pdf.Cell(0, 0, fmt.Sprintf("Volume: %v", r.IssueVol))

	var issueDate string = "Date: -"
	if r.IssueCoverMonth < 12 && r.IssueCoverMonth > 0 {
		month := fmt.Sprintf("%v", time.Month(int(r.IssueCoverMonth)))
		if r.IssueCoverYear > 1 {
			if r.IssueCoverYear == 2 {
				issueDate = "Date: 1899 or before"
			} else {
				issueDate = fmt.Sprintf("Date: %v %v", month, int(r.IssueCoverYear))
			}
		} else { // No cover date year.
			// Do nothing
		}
	} else {
		// No cover year date.
		// Do nothing.
	}
	pdf.SetXY(60, 66)
	pdf.Cell(0, 0, issueDate)

	////
	//// RIGHT SIDE
	////

	pdf.SetFont("Helvetica", "", 10)

	pdf.SetXY(115, 59)
	switch r.SpecialDetails {
	case s_d.SpecialDetailsOther:
		pdf.Cell(0, 0, r.SpecialDetailsOther)
	case s_d.SpecialDetailsRegularEdition:
		pdf.Cell(0, 0, "Regular Edition")
	case s_d.SpecialDetailsDirectEdition:
		pdf.Cell(0, 0, "Direct Edition")
	case s_d.SpecialDetailsNewsstandEdition:
		pdf.Cell(0, 0, "Newstand Edition")
	case s_d.SpecialDetailsVariantCover:
		pdf.Cell(0, 0, "Variant Cover")
	case s_d.SpecialDetailsCanadianPriceVariant:
		pdf.Cell(0, 0, "Canadian Price Variant")
	case s_d.SpecialDetailsFacsimile:
		pdf.Cell(0, 0, "Facsimile")
	case s_d.SpecialDetailsReprint:
		pdf.Cell(0, 0, "Reprint")
	default:
		return nil, fmt.Errorf("missing value for crease finding with %v", r.CreasesFinding)
	}

	title := fmt.Sprintf("%v %v %v", r.SeriesTitle, r.IssueVol, r.IssueNo)
	pdf.SetXY(115, 65)
	pdf.Cell(0, 0, title)

	//
	// GRADING
	//

	switch r.GradingScale {
	case s_d.GradingScaleCPSPercentage:
		pdf.SetFont("Helvetica", "", 24)
		if r.CpsPercentageGrade <= 9 {
			pdf.SetXY(29, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
		} else if r.CpsPercentageGrade <= 99 && r.CpsPercentageGrade > 9 {
			pdf.SetXY(27, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
		} else {
			pdf.SetXY(24, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
		}
	case s_d.GradingScaleNumber:
		pdf.SetFont("Helvetica", "", 60)
		if r.OverallNumberGrade == 10 {
			pdf.SetXY(21.5, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v", r.OverallNumberGrade))
		} else {
			pdf.SetXY(28, 59)
			pdf.Cell(0, 0, fmt.Sprintf("%v", r.OverallNumberGrade))
		}
	case s_d.GradingScaleLetter:
		// If user has chosen the "NM+" option then run the following...
		if r.IsOverallLetterGradeNearMintPlus {
			pdf.SetFont("Helvetica", "", 30)
			pdf.SetXY(26, 57)
			pdf.Cell(0, 0, "NM")

			pdf.SetFont("Helvetica", "", 10)
			pdf.SetXY(22.5, 65)
			pdf.Cell(0, 0, "Near Mint Plus")

			pdf.SetFont("Helvetica", "B", 22) // Start subscript.
			pdf.SetXY(41, 50)
			pdf.Cell(0, 0, "+")
		} else {
			pdf.SetFont("Helvetica", "", 30)
			pdf.SetXY(27, 56)
			pdf.Cell(0, 0, strings.ToUpper(r.OverallLetterGrade))

			pdf.SetFont("Helvetica", "", 14)
			switch r.OverallLetterGrade {
			case "pr", "PR":
				fallthrough
			case "fr", "FR":
				fallthrough
			case "fn", "FN":
				fallthrough
			case "gd", "GD":
				// CASE 1 OF 2: One word description. (Ex: "Fine")
				pdf.SetXY(29, 65)
				pdf.Cell(0, 0, constants.SubmissionOverallLetterGrades[r.OverallLetterGrade])
			case "vg", "VG":
				fallthrough
			case "vf", "VF":
				fallthrough
			case "nm", "NM":
				// CASE 1 OF 2: Two word description. (Ex: "Very Fine")
				pdf.SetXY(23, 65)
				pdf.Cell(0, 0, constants.SubmissionOverallLetterGrades[r.OverallLetterGrade])
			}
		}

		// case s_d.GradingScaleNumber:
		// 	pdf.SetXY(171, 153.5)
		// 	pdf.Cell(0, 0, fmt.Sprintf("%v", r.OverallNumberGrade))
		// case s_d.GradingScaleCPSPercentage:
		// 	pdf.SetXY(171, 153.5)
		// 	pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
	}

	pdf.SetFont("Helvetica", "", 10) // Set back the previous font.

	////
	//// Generate the file and save it to the file.
	////

	fileName := fmt.Sprintf("%s.pdf", r.CPSRN)
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

	return &PDFBuilderResponseDTO{
		FileName: fileName,
		FilePath: filePath,
		Content:  bin,
	}, err
}
