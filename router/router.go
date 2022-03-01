package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-sframe/server"
	"github.com/greatfocus/gf-user/handler"
	"github.com/greatfocus/gf-user/services"
)

// LoadRouter is exported and used in main.go
func LoadRouter(s *server.Meta) *http.ServeMux {
	mux := http.NewServeMux()
	loadHandlers(mux, s)
	log.Println("Created routes with handler")
	return mux
}

// notifyRoute created all routes and handlers relating to controller
func loadHandlers(mux *http.ServeMux, s *server.Meta) {

	// Initialize services
	userService := services.UserService{}
	userService.Init(s)

	otpService := services.OtpService{}
	otpService.Init(s)

	// Initialize routes and handlers
	registerHandler := handler.User{}
	registerHandler.Init(s, &userService)
	mux.Handle("/user/register", server.Use(registerHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.WithoutAuth()))

	forgotPasswordHandler := handler.ForgotPassword{}
	forgotPasswordHandler.Init(s, &userService)
	mux.Handle("/user/forgotpassword", server.Use(forgotPasswordHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.WithoutAuth()))

	otpHandler := handler.Otp{}
	otpHandler.Init(s, &otpService)
	mux.Handle("/user/token", server.Use(otpHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.WithoutAuth()))

	loginHandler := handler.Login{}
	loginHandler.Init(s, &userService)
	mux.Handle("/user/login", server.Use(loginHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.WithoutAuth()))

	userHandler := handler.User{}
	userHandler.Init(s, &userService)
	mux.Handle("/user/users", server.Use(userHandler,
		server.SetHeaders(),
		server.CheckLimitsRates(),
		server.CheckCors(s),
		server.CheckAllowedIPRange(s),
		server.CheckPermission(s),
		server.CheckAuth(s)))
}
