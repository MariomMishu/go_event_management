package cmd

import (
	"ems/controllers"
	"ems/middlewares"
	"ems/routes"
	"ems/server"

	"ems/services"

	"ems/conn"
	db_repo "ems/repositories"

	"ems/config"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: Serve,
}

func Serve(cmd *cobra.Command, args []string) {
	config.LoadConfig()
	// clients
	dbClient := conn.Db()
	redisClient := conn.Redis()

	// repositories
	dbRepo := db_repo.NewRepository(dbClient)

	redisSvc := services.NewRedisService(redisClient)
	userSvc := services.NewUserServiceImpl(dbRepo, redisSvc)
	tokenSvc := services.NewTokenServiceImpl(redisSvc)
	authSvc := services.NewAuthServiceImpl(userSvc, tokenSvc)
	//mailSvc := services.NewMailService(dbRepo)
	userCtrl := controllers.NewUserController(userSvc)
	authCtrl := controllers.NewAuthController(authSvc)
	// middlewares
	authMiddleware := middlewares.NewAuthMiddleware(authSvc, userSvc)

	// Server
	echoServer := echo.New()
	srv := server.New(echoServer)

	routes := routes.New(echoServer, userCtrl, authCtrl, authMiddleware)

	routes.Init()
	srv.Start()

}
