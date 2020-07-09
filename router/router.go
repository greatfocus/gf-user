package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/middlewares"
	"github.com/greatfocus/gf-user/controllers"
)

// Router is exported and used in main.go
func Router(db *database.DB) *http.ServeMux {
	// create new router
	mux := http.NewServeMux()

	// users
	usersRoute(mux, db)

	log.Println("Created routes with controllers")
	return mux
}

// usersRoute created all routes and handlers relating to user controller
func usersRoute(mux *http.ServeMux, db *database.DB) {
	// Initialize controller

	otpController := controllers.OtpController{}
	otpController.Init(db)
	loginController := controllers.LoginController{}
	loginController.Init(db)
	forgotPasswordController := controllers.ForgotPasswordController{}
	forgotPasswordController.Init(db)
	userController := controllers.UserController{}
	userController.Init(db)

	// Initialize routes
	http.HandleFunc("/user/register", middlewares.SetMiddlewareJSON(userController.Handler))
	http.HandleFunc("/user/token", middlewares.SetMiddlewareJSON(otpController.Handler))
	http.HandleFunc("/user/forgotpassword", middlewares.SetMiddlewareJSON(forgotPasswordController.Handler))
	http.HandleFunc("/user/login", middlewares.SetMiddlewareJSON(loginController.Handler))
	http.HandleFunc("/user/users", middlewares.SetMiddlewareJwt(userController.Handler))
}
