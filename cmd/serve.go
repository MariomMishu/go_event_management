package cmd

import (
	"ems/config"
	"ems/conn"
	"ems/controllers"
	"ems/middlewares"
	db_repo "ems/repositories"
	"ems/routes"
	"ems/server"
	"ems/services"
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
	emailClient := conn.EmailClient()
	//worker
	workerPool := conn.WorkerPool()
	// repositories
	dbRepo := db_repo.NewRepository(dbClient)

	redisSvc := services.NewRedisService(redisClient)
	userSvc := services.NewUserServiceImpl(dbRepo, redisSvc)
	tokenSvc := services.NewTokenServiceImpl(redisSvc)
	authSvc := services.NewAuthServiceImpl(userSvc, tokenSvc)
	mailSvc := services.NewMailService(dbRepo, emailClient, workerPool)
	campaignSvc := services.NewCampaignServiceImpl(dbRepo, mailSvc)

	userCtrl := controllers.NewUserController(userSvc)
	authCtrl := controllers.NewAuthController(authSvc)
	campaignCtrl := controllers.NewCampaignController(campaignSvc)
	// middlewares
	authMiddleware := middlewares.NewAuthMiddleware(authSvc, userSvc)

	// Server
	echoServer := echo.New()
	srv := server.New(echoServer)

	routes := routes.New(echoServer, userCtrl, authCtrl, campaignCtrl, authMiddleware)

	routes.Init()
	srv.Start(workerPool)

}
