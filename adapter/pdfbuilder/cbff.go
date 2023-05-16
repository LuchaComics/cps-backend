package pdfbuilder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
)

type CBFFBuilderRequestDTO struct {
	Id        string `json:"id"`
	Uuid      string `json:"uuid"`
	TenantID  string `json:"tenant_id"`
	OrderId   string `json:"order_id"`
	InvoiceID string `json:"invoice_id"`
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
	gofpdi.UseImportedTemplate(pdf, tpl1, 0, 0, 297, 0)

	pdf.SetFont("Courier", "", 11)

	pdf.SetXY(159, 19)
	// pdf.Cell(0, 0, dto.InvoiceID)
	//
	// pdf.SetFont("Helvetica", "", 11)
	//
	// //
	// // Header
	// //
	//
	// pdf.SetXY(116, 25)
	// pdf.Cell(0, 0, dto.InvoiceDate)
	// pdf.SetXY(136, 32)
	// pdf.Cell(0, 0, dto.AssociateName)
	// pdf.SetXY(136, 39)
	// pdf.Cell(0, 0, dto.AssociateTelephone)
	// pdf.SetXY(46, 46)
	// pdf.Cell(0, 0, dto.ClientName)
	// pdf.SetXY(36, 53)
	// pdf.Cell(0, 0, dto.ClientAddress)
	// pdf.SetXY(33, 60)
	// pdf.Cell(0, 0, dto.ClientTelephone)
	// pdf.SetXY(92, 60)
	// pdf.Cell(0, 0, dto.ClientEmail)
	//
	// //
	// // Line #1
	// //
	//
	// pdf.SetXY(30, 72)
	// pdf.Cell(0, 0, dto.Line01Qty)
	// pdf.SetXY(49, 72)
	// pdf.Cell(0, 0, dto.Line01Desc)
	// pdf.SetXY(136, 72)
	// pdf.Cell(0, 0, dto.Line01Price)
	// pdf.SetXY(169, 72)
	// pdf.Cell(0, 0, dto.Line01Amount)
	//
	// //
	// // Line #2
	// //
	//
	// pdf.SetXY(30, 78)
	// if dto.Line02Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line02Qty)
	// }
	// pdf.SetXY(49, 78)
	// pdf.Cell(0, 0, dto.Line02Desc)
	// pdf.SetXY(136, 78)
	// if dto.Line02Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line02Price)
	// }
	// pdf.SetXY(169, 78)
	// if dto.Line02Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line02Amount)
	// }
	//
	// //
	// // Line #3
	// //
	//
	// pdf.SetXY(30, 83)
	// if dto.Line03Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line03Qty)
	// }
	// pdf.SetXY(49, 83)
	// pdf.Cell(0, 0, dto.Line03Desc)
	// pdf.SetXY(136, 83)
	// if dto.Line03Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line03Price)
	// }
	// pdf.SetXY(169, 83)
	// if dto.Line03Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line03Amount)
	// }
	//
	// //
	// // Line #4
	// //
	//
	// pdf.SetXY(30, 88)
	// if dto.Line04Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line04Qty)
	// }
	// pdf.SetXY(49, 88)
	// pdf.Cell(0, 0, dto.Line04Desc)
	// pdf.SetXY(136, 88)
	// if dto.Line04Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line04Price)
	// }
	// pdf.SetXY(169, 88)
	// if dto.Line04Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line04Amount)
	// }
	//
	// //
	// // Line #5
	// //
	//
	// pdf.SetXY(30, 93)
	// if dto.Line05Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line05Qty)
	// }
	// pdf.SetXY(49, 93)
	// pdf.Cell(0, 0, dto.Line05Desc)
	// pdf.SetXY(136, 93)
	// if dto.Line05Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line05Price)
	// }
	// pdf.SetXY(169, 93)
	// if dto.Line05Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line05Amount)
	// }
	//
	// //
	// // Line #6
	// //
	//
	// pdf.SetXY(30, 98)
	// if dto.Line06Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line06Qty)
	// }
	// pdf.SetXY(49, 98)
	// pdf.Cell(0, 0, dto.Line06Desc)
	// pdf.SetXY(136, 98)
	// if dto.Line06Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line06Price)
	// }
	// pdf.SetXY(169, 98)
	// if dto.Line06Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line06Amount)
	// }
	//
	// //
	// // Line #7
	// //
	//
	// pdf.SetXY(30, 103)
	// if dto.Line07Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line07Qty)
	// }
	// pdf.SetXY(49, 103)
	// pdf.Cell(0, 0, dto.Line07Desc)
	// pdf.SetXY(136, 103)
	// if dto.Line07Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line07Price)
	// }
	// pdf.SetXY(169, 103)
	// if dto.Line07Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line07Amount)
	// }
	//
	// //
	// // Line #8
	// //
	//
	// pdf.SetXY(30, 109)
	// if dto.Line07Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line07Amount)
	// }
	// pdf.SetXY(49, 109)
	// pdf.Cell(0, 0, dto.Line08Desc)
	// pdf.SetXY(136, 109)
	// if dto.Line08Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line08Price)
	// }
	// pdf.SetXY(169, 109)
	// if dto.Line08Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line08Amount)
	// }
	//
	// //
	// // Line #9
	// //
	//
	// pdf.SetXY(30, 114)
	// if dto.Line09Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line09Qty)
	// }
	// pdf.SetXY(49, 114)
	// pdf.Cell(0, 0, dto.Line09Desc)
	// pdf.SetXY(136, 114)
	// if dto.Line09Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line09Price)
	// }
	// pdf.SetXY(169, 114)
	// if dto.Line09Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line09Amount)
	// }
	//
	// //
	// // Line #10
	// //
	//
	// pdf.SetXY(30, 119)
	// if dto.Line10Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line10Qty)
	// }
	// pdf.SetXY(49, 119)
	// pdf.Cell(0, 0, dto.Line10Desc)
	// pdf.SetXY(136, 119)
	// if dto.Line10Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line10Price)
	// }
	// pdf.SetXY(169, 119)
	// if dto.Line10Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line10Amount)
	// }
	//
	// //
	// // Line #11
	// //
	//
	// pdf.SetXY(30, 124)
	// if dto.Line11Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line11Qty)
	// }
	// pdf.SetXY(49, 124)
	// pdf.Cell(0, 0, dto.Line11Desc)
	// pdf.SetXY(136, 124)
	// if dto.Line11Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line11Price)
	// }
	// pdf.SetXY(169, 124)
	// if dto.Line11Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line11Amount)
	// }
	//
	// //
	// // Line #12
	// //
	//
	// pdf.SetXY(30, 129)
	// if dto.Line12Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line12Qty)
	// }
	// pdf.SetXY(49, 129)
	// pdf.Cell(0, 0, dto.Line12Desc)
	// pdf.SetXY(136, 129)
	// if dto.Line12Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line12Price)
	// }
	// pdf.SetXY(169, 129)
	// if dto.Line12Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line12Amount)
	// }
	//
	// //
	// // Line #13
	// //
	//
	// pdf.SetXY(30, 134)
	// if dto.Line13Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line13Qty)
	// }
	// pdf.SetXY(49, 134)
	// pdf.Cell(0, 0, dto.Line13Desc)
	// pdf.SetXY(136, 134)
	// if dto.Line13Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line13Price)
	// }
	// pdf.SetXY(169, 134)
	// if dto.Line13Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line13Amount)
	// }
	//
	// //
	// // Line #14
	// //
	//
	// pdf.SetXY(30, 140)
	// if dto.Line14Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line14Qty)
	// }
	// pdf.SetXY(49, 140)
	// pdf.Cell(0, 0, dto.Line14Desc)
	// pdf.SetXY(136, 140)
	// if dto.Line14Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line14Price)
	// }
	// pdf.SetXY(169, 140)
	// if dto.Line14Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line14Amount)
	// }
	//
	// //
	// // Line #15
	// //
	//
	// pdf.SetXY(30, 145)
	// if dto.Line15Qty == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line15Qty)
	// }
	// pdf.SetXY(49, 145)
	// pdf.Cell(0, 0, dto.Line15Desc)
	// pdf.SetXY(136, 145)
	// if dto.Line15Price == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line15Price)
	// }
	// pdf.SetXY(169, 145)
	// if dto.Line15Amount == "0" {
	// 	pdf.Cell(0, 0, "")
	// } else {
	// 	pdf.Cell(0, 0, dto.Line15Amount)
	// }
	//
	// //
	// // Footer
	// //
	//
	// pdf.SetXY(66, 155)
	// pdf.Cell(0, 0, dto.InvoiceQuoteDays)
	// pdf.SetXY(146, 151)
	// pdf.Cell(0, 0, dto.TotalLabour)
	// pdf.SetXY(146, 160)
	// pdf.Cell(0, 0, dto.TotalMaterials)
	// pdf.SetXY(146, 167)
	// pdf.Cell(0, 0, dto.OtherCosts)
	// pdf.SetXY(146, 175)
	// pdf.Cell(0, 0, dto.SubTotal)
	// pdf.SetXY(146, 182)
	// pdf.Cell(0, 0, dto.Tax)
	// pdf.SetXY(146, 192)
	// pdf.Cell(0, 0, dto.Total)
	// pdf.SetXY(77, 170)
	// pdf.Cell(0, 0, dto.InvoiceAssociateTax)
	// pdf.SetXY(77, 179)
	// pdf.Cell(0, 0, dto.InvoiceQuoteDate)
	// pdf.SetXY(77, 187)
	// pdf.Cell(0, 0, dto.InvoiceCustomersApproval)
	// pdf.SetXY(19, 206)
	// pdf.Cell(0, 0, dto.Line01Notes)
	// pdf.SetXY(19, 216)
	// pdf.Cell(0, 0, dto.Line02Notes)
	// pdf.SetXY(55, 225)
	// pdf.Cell(0, 0, dto.PaymentDate)
	// pdf.SetXY(158, 225)
	// pdf.Cell(0, 0, dto.PaymentAmount)
	// pdf.SetXY(158, 217)
	// pdf.Cell(0, 0, dto.Deposit)
	//
	// //
	// // Checks
	// //
	//
	// pdf.SetXY(27, 232)
	// pdf.Cell(0, 0, dto.IsCash)
	// pdf.SetXY(62, 232)
	// pdf.Cell(0, 0, dto.IsCheque)
	// pdf.SetXY(97, 232)
	// pdf.Cell(0, 0, dto.IsDebit)
	// pdf.SetXY(131, 232)
	// pdf.Cell(0, 0, dto.IsCredit)
	// pdf.SetXY(167, 232)
	// pdf.Cell(0, 0, dto.IsOther)
	// pdf.SetXY(42, 248)
	// pdf.Cell(0, 0, dto.AssociateSignDate)
	// pdf.SetXY(155, 248)
	// pdf.Cell(0, 0, dto.AssociateSignature)
	// pdf.SetXY(170, 262)
	// pdf.Cell(0, 0, dto.WorkOrderId)
	//
	// //
	// // Signature
	// //
	//
	// pdf.SetXY(128, 239)
	// pdf.Cell(0, 0, dto.ClientSignature)

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
