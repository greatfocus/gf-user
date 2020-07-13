package router

import (
	"log"
	"net/http"

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
	// Initialize controller

	otpController := controllers.OtpController{}
	otpController.Init(db)
	loginController := controllers.LoginController{}
	loginController.Init(db)
	forgotPasswordController := controllers.ForgotPasswordController{}
	forgotPasswordController.Init(db)
	userController := controllers.UserController{}
	userController.Init(db, config)

	// Initialize routes
	mux.HandleFunc("/user/register", middlewares.SetMiddlewareJSON(userController.Handler))
	mux.HandleFunc("/user/token", middlewares.SetMiddlewareJSON(otpController.Handler))
	mux.HandleFunc("/user/forgotpassword", middlewares.SetMiddlewareJSON(forgotPasswordController.Handler))
	mux.HandleFunc("/user/login", middlewares.SetMiddlewareJSON(loginController.Handler))
	mux.HandleFunc("/user/users", middlewares.SetMiddlewareJwt(userController.Handler))
}
