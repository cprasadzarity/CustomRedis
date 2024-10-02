package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	// custom_redis router
	redisApi := api.Group("/custom_redis")
	SetupRedisRoutes(redisApi)

}
