package main

import (
	"CustomRedis/custom_redis"
	"CustomRedis/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Init routes setup
	routes.Setup(app)

	// Init custom_redis
	custom_redis.Init()

	app.Listen(":8000")
}

// implement delete method
// background to clean up expired ttl
