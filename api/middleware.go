package api

import (
	"github.com/Gedion-01/Go-Crud-Challenge/db"
	"github.com/gofiber/fiber/v2"
)

func CheckUserExists(personStore db.PersonStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		_, found := personStore.Get(userID)
		if !found {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Next()
	}
}
