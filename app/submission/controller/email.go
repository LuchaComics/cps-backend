package controller

import (
	"bytes"
	"context"
	"path"
	"text/template"

	"golang.org/x/exp/slog"
)

func (impl *SubmissionControllerImpl) SendNewSubmissionEmail(email, submissionID, firstName string) error {
	// FOR TESTING PURPOSES ONLY.
	fp := path.Join("templates", "new_submission.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		impl.Logger.Error("parsing error", slog.Any("error", err))
		return err
	}

	var processed bytes.Buffer

	// Render the HTML template with our data.
	data := struct {
		Email      string
		DetailLink string
	}{
		Email:      email,
		DetailLink: "https://" + impl.Emailer.GetDomainName() + "/admin/submission/" + submissionID,
	}
	if err := tmpl.Execute(&processed, data); err != nil {
		impl.Logger.Error("template execution error", slog.Any("error", err))
		return err
	}
	body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

	if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "New Submission", email, body); err != nil {
		impl.Logger.Error("sending error", slog.Any("error", err))
		return err
	}
	return nil
}
