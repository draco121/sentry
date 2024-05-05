package main

import (
	"github.com/draco121/horizon/database"
	"github.com/draco121/horizon/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"sentry/controllers"
	"sentry/core"
	"sentry/repository"
	"sentry/routes"
)

func RunApp() {
	utils.Logger.Info("Starting authorization service")
	client := database.NewMongoDatabase(os.Getenv("MONGODB_URI"))
	db := client.Database("authorization-service")
	repo := repository.NewAuthorizationRepo(db)
	service := core.NewAuthorizationService(client, repo)
	controller := controllers.NewControllers(service)
	router := gin.New()
	router.Use(gin.LoggerWithWriter(utils.Logger.Out))
	routes.RegisterRoutes(controller, router)
	err := router.Run()
	utils.Logger.Info("authorization service started successfully")
	if err != nil {
		utils.Logger.Fatal(err)
		return
	}
}
func main() {
	_ = godotenv.Load()
	RunApp()
}
