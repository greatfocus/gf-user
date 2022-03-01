package router

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/greatfocus/gf-sframe/server"
	"github.com/greatfocus/gf-user/handler"
	"github.com/greatfocus/gf-user/services"
)

// LoadRouter is exported and used in main.go
func LoadRouter(s *server.Server) *http.ServeMux {
	mux := http.NewServeMux()
	loadHandlers(mux, s)
	log.Println("Created routes with handler")
	return mux
}

// notifyRoute created all routes and handlers relating to controller
func loadHandlers(mux *http.ServeMux, s *server.Server) {

	// Initialize services
	userService := services.UserService{}
	userService.Init(s.Database, s.Cache, s.JWT)

	otpService := services.OtpService{}
	otpService.Init(s.Database, s.Cache)

	// Initialize routes and handlers
	registerHandler := handler.User{}
	registerHandler.Init(s, &userService)
	mux.Handle("/user/register", server.Use(registerHandler,
		server.SetHeaders(),
		server.IsThrottle(),
		server.NoAuthentication()))

	forgotPasswordHandler := handler.ForgotPassword{}
	forgotPasswordHandler.Init(s, &userService)
	mux.Handle("/user/forgotpassword", server.Use(forgotPasswordHandler,
		server.SetHeaders(),
		server.IsThrottle(),
		server.NoAuthentication()))

	otpHandler := handler.Otp{}
	otpHandler.Init(s, &otpService)
	mux.Handle("/user/otp", server.Use(otpHandler,
		server.SetHeaders(),
		server.IsThrottle(),
		server.IsAllowedOrigin(os.Getenv("ALLOWED_ORIGIN")),
		server.IsAllowedIPs(os.Getenv("ALLOWED_IP")),
		server.ProcessTimeout(time.Duration(s.Timeout)),
		server.NoAuthentication()))

	loginHandler := handler.Login{}
	loginHandler.Init(s, &userService)
	mux.Handle("/user/login", server.Use(loginHandler,
		server.SetHeaders(),
		server.IsThrottle(),
		server.NoAuthentication()))

	userHandler := handler.User{}
	userHandler.Init(s, &userService)
	mux.Handle("/user/users", server.Use(userHandler,
		server.SetHeaders(),
		server.IsThrottle(),
		server.IsAllowedOrigin(os.Getenv("ALLOWED_ORIGIN")),
		server.IsAllowedIPs(os.Getenv("ALLOWED_IP")),
		server.IsAuthorized(s.JWT),
		server.IsAuthenticated(s.JWT)))
}
