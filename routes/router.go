package routes

import (
	"os"
	"platform_api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func InitRoutes() {
	PORT := os.Getenv("SERVER_PORT")

	challenge := new(controllers.ChallengeController)

	router := gin.Default()

	// recover from panics and respond with internal server error
	router.Use(gin.Recovery())

	// enabling cors
	config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/", challenge.GetAllChallenges)
	router.GET("/id/:id", challenge.GetChallengeById)
	router.GET("/email/:email", challenge.GetChallengeByEmail)

	router.Run(":" + PORT)
}