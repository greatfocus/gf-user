package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-user/services"

	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/middlewares"
	"github.com/greatfocus/gf-user/controllers"
)

// Router is exported and used in main.go
func Router(db *database.DB, config *config.Config) *http.ServeMux {
	// create new router
	mux := http.NewServeMux()

	// users
	usersRoute(mux, db, config)

	log.Println("Created routes with controllers")
	return mux
}

// usersRoute created all routes and handlers relating to user controller
func usersRoute(mux *http.ServeMux, db *database.DB, config *config.Config) {
	// initialize services
	userService := services.UserService{}
	userService.Init(db, config)

	otpService := services.OtpService{}
	otpService.Init(db)

	// Initialize controller
	userController := controllers.UserController{}
	userController.Init(&userService)

	otpController := controllers.OtpController{}
	otpController.Init(&otpService)

	// initialize services
	personService := services.PersonService{}
	personService.Init(db)

	loginController := controllers.LoginController{}
	loginController.Init(&userService)

	forgotPasswordController := controllers.ForgotPasswordController{}
	forgotPasswordController.Init(&userService)

	personController := controllers.PersonController{}
	personController.Init(&personService)

	// Initialize routes
	mux.HandleFunc("/user/register", middlewares.SetMiddlewareJSON(userController.Handler, config.Server.AllowedOrigin))
	mux.HandleFunc("/user/token", middlewares.SetMiddlewareJSON(otpController.Handler, config.Server.AllowedOrigin))
	mux.HandleFunc("/user/forgotpassword", middlewares.SetMiddlewareJSON(forgotPasswordController.Handler, config.Server.AllowedOrigin))
	mux.HandleFunc("/user/login", middlewares.SetMiddlewareJSON(loginController.Handler, config.Server.AllowedOrigin))
	mux.HandleFunc("/user/users", middlewares.SetMiddlewareJwt(userController.Handler, config.Server.AllowedOrigin))
	mux.HandleFunc("/user/person", middlewares.SetMiddlewareJwt(personController.Handler, config.Server.AllowedOrigin))
}
