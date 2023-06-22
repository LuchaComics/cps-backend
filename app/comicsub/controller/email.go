package controller

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"text/template"

	"golang.org/x/exp/slog"

	s_d "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
)

func (impl *ComicSubmissionControllerImpl) sendNewComicSubmissionEmails(m *s_d.ComicSubmission) error {
	//
	// ROOT
	//

	impl.Logger.Debug("sending to root staff",
		slog.Any("submission-id", m.ID))

	response, err := impl.UserStorer.ListAllRootStaff(context.Background())
	if err != nil {
		impl.Logger.Error("database list all staff error", slog.Any("error", err))
		return err
	}

	for _, u := range response.Results {
		if err := impl.sendStaffNewComicSubmissionEmail(u.Email, m); err != nil {
			impl.Logger.Error("failed sending stafff email error", slog.Any("error", err))
			return err
		}
	}

	//
	// RETAILERS
	//

	impl.Logger.Debug("sending to all retailer staff",
		slog.Any("submission-id", m.ID),
		slog.Any("organization-id", m.OrganizationID))

	response, err = impl.UserStorer.ListAllRetailerStaffForOrganizationID(context.Background(), m.OrganizationID)
	if err != nil {
		impl.Logger.Error("database list all retailer error", slog.Any("error", err))
		return err
	}

	for _, u := range response.Results {
		if err := impl.sendRetailerNewComicSubmissionEmail(u.Email, m); err != nil {
			impl.Logger.Error("failed sending retailer email error", slog.Any("error", err))
			return err
		}
	}
	return nil
}

func (impl *ComicSubmissionControllerImpl) sendStaffNewComicSubmissionEmail(staffEmail string, m *s_d.ComicSubmission) error {
	// FOR TESTING PURPOSES ONLY.
	fp := path.Join("templates", "staff_submission_created.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		impl.Logger.Error("parsing error", slog.Any("error", err))
		return err
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		OrganizationName string
		Item             string
		CPSRN            string
		DetailLink       string
	}{
		OrganizationName: m.OrganizationName,
		Item:             m.Item,
		CPSRN:            m.CPSRN,
		DetailLink:       fmt.Sprintf("https://%v/admin/submission/%v", impl.Emailer.GetDomainName(), m.ID.Hex()),
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "New Comic Submission", staffEmail, body); err != nil {
		impl.Logger.Error("sending error", slog.Any("error", err))
		return err
	}
	impl.Logger.Debug("sent `New Comic Submission` email",
		slog.String("staff-email", staffEmail),
		slog.Any("submission-id", m.ID))
	return nil
}

func (impl *ComicSubmissionControllerImpl) sendRetailerNewComicSubmissionEmail(retailerEmail string, m *s_d.ComicSubmission) error {
	// FOR TESTING PURPOSES ONLY.
	fp := path.Join("templates", "retailer_submission_created.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		impl.Logger.Error("parsing error", slog.Any("error", err))
		return err
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		OrganizationName string
		Item             string
		CPSRN            string
		DetailLink       string
	}{
		OrganizationName: m.OrganizationName,
		Item:             m.Item,
		CPSRN:            m.CPSRN,
		DetailLink:       fmt.Sprintf("https://%v/submission/%v", impl.Emailer.GetDomainName(), m.ID.Hex()),
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "Submitted to CPS", retailerEmail, body); err != nil {
		impl.Logger.Error("sending error", slog.Any("error", err))
		return err
	}
	impl.Logger.Debug("sent `Submitted to CPS` email",
		slog.String("retailer-email", retailerEmail),
		slog.Any("submission-id", m.ID),
		slog.Any("organization-id", m.OrganizationID))
	return nil
}
