package http

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/exp/slog"

	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/inputport/http/attachment"
	"github.com/LuchaComics/cps-backend/inputport/http/comicsub"
	"github.com/LuchaComics/cps-backend/inputport/http/customer"
	"github.com/LuchaComics/cps-backend/inputport/http/gateway"
	"github.com/LuchaComics/cps-backend/inputport/http/middleware"
	"github.com/LuchaComics/cps-backend/inputport/http/organization"
	"github.com/LuchaComics/cps-backend/inputport/http/user"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type httpInputPort struct {
	Config          *config.Conf
	Logger          *slog.Logger
	Server          *http.Server
	Middleware      middleware.Middleware
	Gateway         *gateway.Handler
	User            *user.Handler
	Organization    *organization.Handler
	ComicSubmission *comicsub.Handler
	Customer        *customer.Handler
	Attachment      *attachment.Handler
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	mid middleware.Middleware,
	gh *gateway.Handler,
	cu *user.Handler,
	org *organization.Handler,
	t *comicsub.Handler,
	cust *customer.Handler,
	att *attachment.Handler,
) InputPortServer {
	// Initialize the ServeMux.
	mux := http.NewServeMux()

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation via `https://github.com/rs/cors` for more options.
	handler := cors.AllowAll().Handler(mux)

	// Bind the HTTP server to the assigned address and port.
	addr := fmt.Sprintf("%s:%s", configp.AppServer.IP, configp.AppServer.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Create our HTTP server controller.
	p := &httpInputPort{
		Config:          configp,
		Logger:          loggerp,
		Middleware:      mid,
		Gateway:         gh,
		User:            cu,
		Organization:    org,
		ComicSubmission: t,
		Customer:        cust,
		Attachment:      att,
		Server:          srv,
	}

	// Attach the HTTP server controller to the ServerMux.
	mux.HandleFunc("/", mid.Attach(p.HandleRequests))

	return p
}

func (port *httpInputPort) Run() {
	port.Logger.Info("HTTP server running")
	if err := port.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		port.Logger.Error("listen failed", slog.Any("error", err))

		// DEVELOPERS NOTE: We terminate app here b/c dependency injection not allowed to fail, so fail here at startup of app.
		panic("failed running")
	}
}

func (port *httpInputPort) Shutdown() {
	port.Logger.Info("HTTP server shutdown")
}

func (port *httpInputPort) HandleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get our URL paths which are slash-seperated.
	ctx := r.Context()
	p := ctx.Value("url_split").([]string)
	n := len(p)
	port.Logger.Debug("Handling request",
		slog.Int("n", n),
		slog.String("m", r.Method),
		slog.Any("p", p),
	)

	switch {
	// --- GATEWAY & PROFILE & DASHBOARD --- //
	case n == 3 && p[1] == "v1" && p[2] == "version" && r.Method == http.MethodGet:
		port.Gateway.Version(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "greeting" && r.Method == http.MethodPost:
		port.Gateway.Greet(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "login" && r.Method == http.MethodPost:
		port.Gateway.Login(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "register" && r.Method == http.MethodPost:
		port.Gateway.Register(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "refresh-token" && r.Method == http.MethodPost:
		port.Gateway.RefreshToken(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "verify" && r.Method == http.MethodPost:
		port.Gateway.Verify(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "logout" && r.Method == http.MethodPost:
		port.Gateway.Logout(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "profile" && r.Method == http.MethodGet:
		port.Gateway.Profile(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "profile" && r.Method == http.MethodPut:
		port.Gateway.ProfileUpdate(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "profile" && p[3] == "change-password" && r.Method == http.MethodPut:
		port.Gateway.ProfileChangePassword(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "forgot-password" && r.Method == http.MethodPost:
		port.Gateway.ForgotPassword(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "password-reset" && r.Method == http.MethodPost:
		port.Gateway.PasswordReset(w, r)
		// case n == 3 && p[1] == "v1" && p[2] == "profile" && r.Method == http.MethodGet:
		// ...

		// --- REGISTRY --- // (TODO)
	case n == 4 && p[1] == "v1" && p[2] == "cpsrn" && r.Method == http.MethodGet:
		port.ComicSubmission.GetRegistryByCPSRN(w, r, p[3])

		// --- SUBMISSIONS --- //
	case n == 3 && p[1] == "v1" && p[2] == "comic-submissions" && r.Method == http.MethodGet:
		port.ComicSubmission.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "comic-submissions" && r.Method == http.MethodPost:
		port.ComicSubmission.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "comic-submission" && r.Method == http.MethodGet:
		port.ComicSubmission.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "comic-submission" && r.Method == http.MethodPut:
		port.ComicSubmission.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "comic-submission" && r.Method == http.MethodDelete:
		port.ComicSubmission.ArchiveByID(w, r, p[3])
	case n == 5 && p[1] == "v1" && p[2] == "comic-submission" && p[4] == "perma-delete" && r.Method == http.MethodDelete:
		port.ComicSubmission.DeleteByID(w, r, p[3])
	case n == 5 && p[1] == "v1" && p[2] == "comic-submissions" && p[3] == "operation" && p[4] == "set-user" && r.Method == http.MethodPost:
		port.ComicSubmission.OperationSetUser(w, r)
	case n == 5 && p[1] == "v1" && p[2] == "comic-submissions" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
		port.ComicSubmission.OperationCreateComment(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "comic-submissions" && p[3] == "select-options" && r.Method == http.MethodGet:
		port.ComicSubmission.ListAsSelectOptionByFilter(w, r)

	// --- ORGANIZATION --- //
	case n == 3 && p[1] == "v1" && p[2] == "organizations" && r.Method == http.MethodGet:
		port.Organization.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "organizations" && r.Method == http.MethodPost:
		port.Organization.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "organization" && r.Method == http.MethodGet:
		port.Organization.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "organization" && r.Method == http.MethodPut:
		port.Organization.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "organization" && r.Method == http.MethodDelete:
		port.Organization.DeleteByID(w, r, p[3])
	case n == 5 && p[1] == "v1" && p[2] == "organizations" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
		port.Organization.OperationCreateComment(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "organizations" && p[3] == "select-options" && r.Method == http.MethodGet:
		port.Organization.ListAsSelectOptionByFilter(w, r)

	// --- CUSTOMERS --- //
	case n == 3 && p[1] == "v1" && p[2] == "customers" && r.Method == http.MethodGet:
		port.Customer.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "customers" && r.Method == http.MethodPost:
		port.Customer.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "customer" && r.Method == http.MethodGet:
		port.Customer.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "customer" && r.Method == http.MethodPut:
		port.Customer.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "customer" && r.Method == http.MethodDelete:
		port.Customer.DeleteByID(w, r, p[3])
	case n == 5 && p[1] == "v1" && p[2] == "customers" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
		port.Customer.OperationCreateComment(w, r)

	// --- USERS --- //
	case n == 3 && p[1] == "v1" && p[2] == "users" && r.Method == http.MethodGet:
		port.User.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "users" && r.Method == http.MethodPost:
		port.User.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "user" && r.Method == http.MethodGet:
		port.User.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "user" && r.Method == http.MethodPut:
		port.User.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "user" && r.Method == http.MethodDelete:
		port.User.DeleteByID(w, r, p[3])
	case n == 5 && p[1] == "v1" && p[2] == "users" && p[3] == "operation" && p[4] == "create-comment" && r.Method == http.MethodPost:
		port.User.OperationCreateComment(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "users" && p[3] == "select-options" && r.Method == http.MethodGet:
		port.User.ListAsSelectOptions(w, r)

	// --- ATTACHMENTS --- //
	case n == 3 && p[1] == "v1" && p[2] == "attachments" && r.Method == http.MethodGet:
		port.Attachment.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "attachments" && r.Method == http.MethodPost:
		port.Attachment.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodGet:
		port.Attachment.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodPut:
		port.Attachment.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "attachment" && r.Method == http.MethodDelete:
		port.Attachment.DeleteByID(w, r, p[3])

	// --- CATCH ALL: D.N.E. ---
	default:
		http.NotFound(w, r)
	}
}
