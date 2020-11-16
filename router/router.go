package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/gf-user/services"

	"github.com/greatfocus/gf-frame/middlewares"
	"github.com/greatfocus/gf-frame/server"
	"github.com/greatfocus/gf-user/controllers"
)

// Router is exported and used in main.go
func Router(s *server.Server) *http.ServeMux {
	// create new router
	mux := http.NewServeMux()

	// users
	usersRoute(mux, s)

	log.Println("Created routes with controllers")
	return mux
}

// usersRoute created all routes and handlers relating to user controller
func usersRoute(mux *http.ServeMux, server *server.Server) {
	// initialize services
	userService := services.UserService{}
	userService.Init(server)

	otpService := services.OtpService{}
	otpService.Init(server)

	personService := services.PersonService{}
	personService.Init(server)

	clientService := services.ClientService{}
	clientService.Init(server)

	// Initialize controller
	userController := controllers.UserController{}
	userController.Init(&userService)

	otpController := controllers.OtpController{}
	otpController.Init(&otpService)

	loginController := controllers.LoginController{}
	loginController.Init(&userService)

	forgotPasswordController := controllers.ForgotPasswordController{}
	forgotPasswordController.Init(&userService)

	personController := controllers.PersonController{}
	personController.Init(&personService)

	contactController := controllers.ContactController{}
	contactController.Init(&userService)

	clientController := controllers.ClientController{}
	clientController.Init(&clientService)

	clientAuthController := controllers.ClientAuthController{}
	clientAuthController.Init(&clientService)

	// Initialize routes
	mux.HandleFunc("/user/register", middlewares.SetMiddlewareJSON(userController.Handler, server))
	mux.HandleFunc("/user/token", middlewares.SetMiddlewareJSON(otpController.Handler, server))
	mux.HandleFunc("/user/forgotpassword", middlewares.SetMiddlewareJSON(forgotPasswordController.Handler, server))
	mux.HandleFunc("/user/login", middlewares.SetMiddlewareJSON(loginController.Handler, server))
	mux.HandleFunc("/user/users", middlewares.SetMiddlewareJwt(userController.Handler, server))
	mux.HandleFunc("/user/person", middlewares.SetMiddlewareJwt(personController.Handler, server))
	mux.HandleFunc("/user/contact", middlewares.SetMiddlewareJSON(contactController.Handler, server))
	mux.HandleFunc("/user/client/login", middlewares.SetMiddlewareJSON(clientAuthController.Handler, server))
	mux.HandleFunc("/user/client", middlewares.SetMiddlewareJwt(clientController.Handler, server))

}
