package http

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/exp/slog"

	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/inputport/http/gateway"
	"github.com/LuchaComics/cps-backend/inputport/http/middleware"
	"github.com/LuchaComics/cps-backend/inputport/http/submission"
	"github.com/LuchaComics/cps-backend/inputport/http/user"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type httpInputPort struct {
	Config     *config.Conf
	Logger     *slog.Logger
	Server     *http.Server
	Middleware middleware.Middleware
	Gateway    *gateway.Handler
	User       *user.Handler
	Submission *submission.Handler
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	mid middleware.Middleware,
	gh *gateway.Handler,
	cu *user.Handler,
	t *submission.Handler,
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
		Config:     configp,
		Logger:     loggerp,
		Middleware: mid,
		Gateway:    gh,
		User:       cu,
		Submission: t,
		Server:     srv,
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
		// case n == 3 && p[1] == "v1" && p[2] == "profile" && r.Method == http.MethodGet:
		// ...

		// --- SUBMISSIONS --- //
	case n == 3 && p[1] == "v1" && p[2] == "submissions" && r.Method == http.MethodGet:
		port.Submission.List(w, r)
	case n == 3 && p[1] == "v1" && p[2] == "submissions" && r.Method == http.MethodPost:
		port.Submission.Create(w, r)
	case n == 4 && p[1] == "v1" && p[2] == "submission" && r.Method == http.MethodGet:
		port.Submission.GetByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "submission" && r.Method == http.MethodPut:
		port.Submission.UpdateByID(w, r, p[3])
	case n == 4 && p[1] == "v1" && p[2] == "submission" && r.Method == http.MethodDelete:
		port.Submission.DeleteByID(w, r, p[3])

	// --- CATCH ALL: D.N.E. ---
	default:
		http.NotFound(w, r)
	}
}
