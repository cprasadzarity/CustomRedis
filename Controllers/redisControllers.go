package Controllers

import (
	"CustomRedis/common"
	"CustomRedis/custom_redis"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func GetController(ctx *fiber.Ctx) error {
	key := ctx.Params("key")

	value, err := custom_redis.Rds.Get(key)
	if err != nil {
		if customError, ok := err.(*common.CustomError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": customError,
			})
		} else {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    value,
		"message": nil,
	})
}

func SetController(ctx *fiber.Ctx) error {
	var data map[string]string
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	err = common.ValidateRequiredKeys(data, "key", "value", "ttl")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	ttl_seconds, err := strconv.Atoi(data["ttl"])
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid ttl",
		})
	}

	key := data["key"]
	value := data["value"]
	ttl := time.Duration(ttl_seconds) * time.Second

	err = custom_redis.Rds.Set(key, value, ttl)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    value,
	})
}
