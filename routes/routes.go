package routes

import (
	"ems/consts"
	"ems/controllers"
	m "ems/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Route struct {
	echo           *echo.Echo
	userCtrl       *controllers.UserController
	authCtrl       *controllers.AuthController
	authMiddleware *m.AuthMiddleware
}

func New(e *echo.Echo, userCtrl *controllers.UserController, authCtrl *controllers.AuthController, authMiddleware *m.AuthMiddleware) *Route {
	return &Route{
		echo:           e,
		userCtrl:       userCtrl,
		authCtrl:       authCtrl,
		authMiddleware: authMiddleware,
	}
}

func (r *Route) Init() {
	e := r.echo
	m.Init(e) //Middleware Initialization
	// APM routes
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	g := e.Group("/v1")
	user := g.Group("/user")
	user.POST("/create", r.userCtrl.SignUp)
	// user.POST("/profile", r.userCtrl.GetProfile, m.Authenticate("user"))
	user.GET("/profile", r.userCtrl.GetProfile, r.authMiddleware.Authenticate(""))
	user.POST("", r.userCtrl.CreateUser, r.authMiddleware.Authenticate(consts.PermissionUserCreate))
	auth := g.Group("/auth")
	auth.POST("/login", r.authCtrl.Login)
	auth.POST("/logout", r.authCtrl.Logout, r.authMiddleware.Authenticate(""))
}
