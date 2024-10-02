package routes

import (
	"CustomRedis/Controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRedisRoutes(router fiber.Router) {
	router.Get("/:key", Controllers.GetController)
	router.Post("/", Controllers.SetController)
}
