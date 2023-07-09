package pdfbuilder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
	"golang.org/x/exp/slog"

	s_d "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	c "github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// CPS PEDIGREE COLLECTION

type PCBuilderRequestDTO struct {
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
	PrimaryLabelDetails                     int8                       `bson:"primary_label_details" json:"primary_label_details"`
	PrimaryLabelDetailsOther                string                     `bson:"primary_label_details_other" json:"primary_label_details_other"`
}

type PCBuilder interface {
	GeneratePDF(dto *PCBuilderRequestDTO) (*PDFBuilderResponseDTO, error)
}

type pcBuilder struct {
	PDFTemplateFilePath string
	DataDirectoryPath   string
	UUID                uuid.Provider
	Logger              *slog.Logger
}

func NewPCBuilder(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) PCBuilder {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("pdf builder for pc initializing...")
	_, err := os.Stat(cfg.PDFBuilder.PCTemplatePath)
	if os.IsNotExist(err) {
		log.Fatal(errors.New("file does not exist"))
	}

	return &pcBuilder{
		PDFTemplateFilePath: cfg.PDFBuilder.PCTemplatePath,
		DataDirectoryPath:   cfg.PDFBuilder.DataDirectoryPath,
		UUID:                uuidp,
		Logger:              logger,
	}
}

func (bdr *pcBuilder) GeneratePDF(r *PCBuilderRequestDTO) (*PDFBuilderResponseDTO, error) {
	specialNotesLines := splitText(r.SpecialNotes, 50)

	var err error

	// Open our PDF invoice template and create clone it for the PDF invoice we will be building with.
	pdf := gofpdf.New("P", "mm", "A3", "")
	tpl1 := gofpdi.ImportPage(pdf, bdr.PDFTemplateFilePath, 1, "/MediaBox")

	pdf.AddPage()

	// Draw imported template onto page
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 0, 297, 420)

	/////
	///// TOP PAEG
	/////

	pdf.SetFont("Helvetica", "B", 12)

	// ROW 1
	pdf.SetXY(74+65, 24.5)
	pdf.Cell(0, 0, r.SeriesTitle)

	// ROW 2
	pdf.SetXY(80, 32.5)
	pdf.Cell(0, 0, r.IssueVol)

	pdf.SetXY(126, 32.5)
	pdf.Cell(0, 0, r.IssueNo)

	pdf.SetXY(126, 32.5)
	pdf.Cell(0, 0, r.IssueNo)

	if r.IssueCoverMonth < 12 && r.IssueCoverMonth > 0 {
		pdf.SetXY(187, 32.5)
		pdf.Cell(0, 0, fmt.Sprintf("%v", time.Month(int(r.IssueCoverMonth))))
	} else {
		pdf.SetXY(187, 32.5)
		pdf.Cell(0, 0, "-") // No cover year date.
	}
	if r.IssueCoverYear > 1 {
		pdf.SetXY(205, 32.5)
		if r.IssueCoverYear == 2 {
			pdf.Cell(0, 0, "1899 or before")
		} else {
			pdf.Cell(0, 0, fmt.Sprintf("%v", int(r.IssueCoverYear)))
		}
	} else { // No cover date year.
		pdf.SetXY(187, 32.5)
		pdf.Cell(0, 0, "-")
	}

	// ROW 3
	pdf.SetXY(100, 40)
	if r.PrimaryLabelDetails == 1 {
		pdf.Cell(0, 0, r.PrimaryLabelDetailsOther)
	} else {
		pdf.Cell(0, 0, constants.SubmissionPrimaryLabelDetails[r.PrimaryLabelDetails])
	}

	pdf.SetFont("Helvetica", "B", 55)

	// ROW 10 - Grading
	switch r.GradingScale {
	case s_d.GradingScaleLetter:
		pdf.SetXY(246, 30)
		pdf.Cell(0, 0, strings.ToUpper(r.OverallLetterGrade))

		// If user has chosen the "NM+" option then run the following...
		if r.IsOverallLetterGradeNearMintPlus {
			pdf.SetFont("Helvetica", "B", 25) // Start subscript.
			pdf.SetXY(276, 23)
			pdf.Cell(0, 0, "+")
			pdf.SetFont("Helvetica", "B", 40) // Resume the previous font.
		}
	case s_d.GradingScaleNumber:
		pdf.SetXY(243, 30)
		pdf.Cell(0, 0, fmt.Sprintf("%v", r.OverallNumberGrade))
	case s_d.GradingScaleCPSPercentage:
		pdf.SetXY(243, 30)
		pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
	}

	/////
	///// BOTTOM PAEG
	/////

	//
	// CPS REGISTRY NUMBER
	//

	pdf.SetFont("Courier", "", 12)
	pdf.SetXY(17, 305)
	pdf.Cell(0, 0, r.CPSRN)

	//
	// LEFT SIDE
	//

	pdf.SetFont("Helvetica", "B", 8)

	// ROW 1
	pdf.SetXY(77, 320)
	pdf.Cell(0, 0, fmt.Sprintf("%v", r.SubmissionDate.Day())) // Day
	pdf.SetXY(87, 320)
	pdf.Cell(0, 0, fmt.Sprintf("%v", int(r.SubmissionDate.Month()))) // Month (number)
	pdf.SetXY(93, 320)
	pdf.Cell(0, 0, fmt.Sprintf("%v", r.SubmissionDate.Year())) // Day
	//
	// ROW 2
	pdf.SetXY(57, 325)
	pdf.Cell(0, 0, r.UserFirstName)
	pdf.SetXY(80, 325)
	pdf.Cell(0, 0, r.UserLastName)

	// ROW 3
	pdf.SetXY(40, 330)
	pdf.Cell(0, 0, r.UserOrganizationName)

	//
	// RIGHT SIDE
	//

	// ROW 1
	pdf.SetXY(140, 320)
	pdf.Cell(0, 0, r.SeriesTitle)

	// ROW 2
	pdf.SetXY(112, 325)
	pdf.Cell(0, 0, r.IssueVol)
	pdf.SetXY(135, 325)
	pdf.Cell(0, 0, r.IssueNo)
	if r.IssueCoverMonth < 12 && r.IssueCoverMonth > 0 {
		pdf.SetXY(163, 325)
		pdf.Cell(0, 0, fmt.Sprintf("%v", time.Month(int(r.IssueCoverMonth))))
	} else {
		pdf.SetXY(163, 325)
		pdf.Cell(0, 0, "-") // No cover year date.
	}
	if r.IssueCoverYear > 1 {
		pdf.SetXY(178, 325)
		if r.IssueCoverYear == 2 {
			pdf.Cell(0, 0, "1899 or before")
		} else {
			pdf.Cell(0, 0, fmt.Sprintf("%v", int(r.IssueCoverYear)))
		}
	} else { // No cover date year.
		pdf.SetXY(257, 47.5)
		pdf.Cell(0, 0, "-")
	}

	// ROW 3
	pdf.SetXY(142, 330)
	pdf.Cell(0, 0, r.PublisherName)

	//
	// RIGHT
	//

	pdf.SetFont("Helvetica", "B", 8) // This controls the next text.

	// ROW 1 - Creases
	switch strings.ToLower(r.CreasesFinding) {
	case "pr":
		pdf.SetXY(63, 341.5)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(76, 341.5)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(88, 341.5)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(100, 341.5)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(113, 341.5)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(125, 341.5)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(137, 341.5)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, fmt.Errorf("missing value for crease finding with %v", r.CreasesFinding)
	}

	// ROW 2 - Tears
	switch strings.ToLower(r.TearsFinding) {
	case "pr":
		pdf.SetXY(63, 346)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(76, 346)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(88, 346)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(100, 346)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(113, 346)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(125, 346)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(137, 346)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for tears finding")
	}

	// ROW 3 - Missing Parts
	switch strings.ToLower(r.MissingPartsFinding) {
	case "pr":
		pdf.SetXY(63, 351)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(76, 351)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(88, 351)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(100, 351)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(113, 351)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(125, 351)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(137, 351)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for missing parts finding")
	}

	// ROW 4 - Stains / Marks / Substances
	switch strings.ToLower(r.StainsFinding) {
	case "pr":
		pdf.SetXY(63, 355)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(76, 355)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(88, 355)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(100, 355)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(113, 355)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(125, 355)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(137, 355)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for stains finding")
	}

	// ROW 5 - Distortion / Colour
	switch strings.ToLower(r.DistortionFinding) {
	case "pr":
		pdf.SetXY(63, 359.5)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(76, 359.5)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(88, 359.5)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(100, 359.5)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(113, 359.5)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(125, 359.5)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(137, 359.5)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for distorion finding")
	}

	// ROW 6 - Paper Quality
	switch strings.ToLower(r.PaperQualityFinding) {
	case "pr":
		pdf.SetXY(63, 364)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(76, 364)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(88, 364)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(100, 364)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(113, 364)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(125, 364)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(137, 364)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for paper quality finding")
	}

	// ROW 7 - Spine / Staples
	switch strings.ToLower(r.SpineFinding) {
	case "pr":
		pdf.SetXY(63, 369)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(76, 369)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(88, 369)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(100, 369)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(113, 369)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(125, 369)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(137, 369)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value for paper quality finding")
	}

	// ROW 8 - Cover (Front & Back)
	switch strings.ToLower(r.CoverFinding) {
	case "pr":
		pdf.SetXY(63, 373)
		pdf.Cell(0, 0, "PR")
	case "fr":
		pdf.SetXY(76, 373)
		pdf.Cell(0, 0, "FR")
	case "gd":
		pdf.SetXY(88, 373)
		pdf.Cell(0, 0, "GD")
	case "vg":
		pdf.SetXY(100, 373)
		pdf.Cell(0, 0, "VG")
	case "fn":
		pdf.SetXY(113, 373)
		pdf.Cell(0, 0, "FN")
	case "vf":
		pdf.SetXY(125, 373)
		pdf.Cell(0, 0, "VF")
	case "nm":
		pdf.SetXY(137, 373)
		pdf.Cell(0, 0, "NM")
	default:
		return nil, errors.New("missing value cover finding")
	}

	// ROW 9 - Shows signs of temp
	if r.ShowsSignsOfTamperingOrRestoration == true {
		pdf.SetXY(59, 377.5)
		pdf.Cell(0, 0, "X")
	} else {
		pdf.SetXY(69.5, 377.5)
		pdf.Cell(0, 0, "X")
	}

	pdf.SetFont("Helvetica", "B", 30)

	// ROW 10 - Grading
	switch r.GradingScale {
	case s_d.GradingScaleLetter:
		pdf.SetXY(117, 388)
		pdf.Cell(0, 0, strings.ToUpper(r.OverallLetterGrade))

		// If user has chosen the "NM+" option then run the following...
		if r.IsOverallLetterGradeNearMintPlus {
			pdf.SetFont("Helvetica", "B", 20) // Start subscript.
			pdf.SetXY(133, 385)
			pdf.Cell(0, 0, "+")
			pdf.SetFont("Helvetica", "B", 40) // Resume the previous font.
		}
	case s_d.GradingScaleNumber:
		pdf.SetXY(117, 388)
		pdf.Cell(0, 0, fmt.Sprintf("%v", r.OverallNumberGrade))
	case s_d.GradingScaleCPSPercentage:
		pdf.SetXY(117, 388)
		pdf.Cell(0, 0, fmt.Sprintf("%v%%", r.CpsPercentageGrade))
	}

	//
	// LEFT
	//

	pdf.SetFont("Helvetica", "", 5)

	if len(r.SpecialNotes) > 638 {
		return nil, errors.New("special notes length over 455")
	}

	if specialNote, ok := getElementAtIndex(specialNotesLines, 0); ok { // ROW 1
		pdf.SetXY(150, 339+1.85*0)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 1); ok {
		pdf.SetXY(150, 339+1.85*1)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 2); ok {
		pdf.SetXY(150, 339+1.85*2)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 3); ok {
		pdf.SetXY(150, 339+1.85*3)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 4); ok {
		pdf.SetXY(150, 339+1.85*4)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 5); ok {
		pdf.SetXY(150, 339+1.85*5)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 6); ok {
		pdf.SetXY(150, 339+1.85*6)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 7); ok {
		pdf.SetXY(150, 339+1.85*7)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 8); ok {
		pdf.SetXY(150, 339+1.85*8)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 9); ok {
		pdf.SetXY(150, 339+1.85*9)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 10); ok {
		pdf.SetXY(150, 339+1.85*10)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 11); ok {
		pdf.SetXY(150, 339+1.85*11)
		pdf.Cell(0, 0, specialNote)
	}
	if specialNote, ok := getElementAtIndex(specialNotesLines, 12); ok { // ROW 13 - MAXIMUM
		pdf.SetXY(150, 339+1.85*12)
		pdf.Cell(0, 0, specialNote)
	}

	////////////////////////////////////////////////////////////////////////////

	if len(r.GradingNotes) > 638 {
		return nil, errors.New("grading notes length over 638")
	}

	gradingNotesLines := splitText(r.GradingNotes, 50)

	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 0); ok { // ROW 1
		pdf.SetXY(150, 369+1.85*0)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 1); ok {
		pdf.SetXY(150, 369+1.85*1)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 2); ok {
		pdf.SetXY(150, 369+1.85*2)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 3); ok {
		pdf.SetXY(150, 369+1.85*3)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 4); ok {
		pdf.SetXY(150, 369+1.85*4)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 5); ok {
		pdf.SetXY(150, 369+1.85*5)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 6); ok {
		pdf.SetXY(150, 369+1.85*6)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 7); ok {
		pdf.SetXY(150, 369+1.85*7)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 8); ok {
		pdf.SetXY(150, 369+1.85*8)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 9); ok {
		pdf.SetXY(150, 369+1.85*9)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 10); ok {
		pdf.SetXY(150, 369+1.85*10)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 11); ok {
		pdf.SetXY(150, 369+1.85*11)
		pdf.Cell(0, 0, gradingNote)
	}
	if gradingNote, ok := getElementAtIndex(gradingNotesLines, 12); ok {
		pdf.SetXY(150, 369+1.85*12)
		pdf.Cell(0, 0, gradingNote)
	}

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
