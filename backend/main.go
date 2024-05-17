package main

import (
	_ "fmt"

	"log"

	"github.com/gooddavvy/markdown-repository-snapshot-app/backend/utils"

	_ "github.com/joho/godotenv/autoload" // Load the .env file

	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	app = fiber.New() // Create the backend server
)

func main() {
	app.Use(logger.New())
	viper.AutomaticEnv()

	app.Post("/gen-md", utils.GenMdHandler)
	log.Fatal(app.Listen(":8081"))

}
