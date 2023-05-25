package mailgun

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

type Emailer interface {
	Send(ctx context.Context, sender, subject, recipient, htmlContent string) error
	GetSenderEmail() string
	GetDomainName() string
}

type mailgunEmailer struct {
	Mailgun     *mailgun.MailgunImpl
	UUID        uuid.Provider
	Logger      *slog.Logger
	senderEmail string
	domainName  string
}

func NewEmailer(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) Emailer {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("mailgun emailer initializing...")
	mg := mailgun.NewMailgun(cfg.Emailer.Domain, cfg.Emailer.APIKey)
	logger.Debug("mailgun emailer was initialized.")

	mg.SetAPIBase(cfg.Emailer.APIBase) // Override to support our custom email requirements.

	return &mailgunEmailer{
		Mailgun:     mg,
		UUID:        uuidp,
		Logger:      logger,
		senderEmail: cfg.Emailer.SenderEmail,
		domainName:  cfg.AppServer.DomainName,
	}
}

func (me *mailgunEmailer) Send(ctx context.Context, sender, subject, recipient, body string) error {
	me.Logger.Debug("sent email",
		slog.String("sender", sender),
		slog.String("subject", subject),
		slog.String("recipient", recipient),
		slog.String("body", body))

	message := me.Mailgun.NewMessage(sender, subject, "", recipient)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := me.Mailgun.Send(ctx, message)

	if err != nil {
		me.Logger.Error("emailer failed sending", slog.Any("err", err))
		return err
	}

	me.Logger.Debug("emailer sent with response", slog.Any("id", id), slog.Any("resp", resp))

	return nil
}

func (me *mailgunEmailer) GetSenderEmail() string {
	return me.senderEmail
}

func (me *mailgunEmailer) GetDomainName() string {
	return me.domainName
}
