package api

import (
	"github.com/Gedion-01/Go-Crud-Challenge/db"
	"github.com/Gedion-01/Go-Crud-Challenge/types"
	"github.com/gofiber/fiber/v2"
)

type PersonHandler struct {
	personStore db.PersonStore
}

func NewPersonHandler(personStore db.PersonStore) *PersonHandler {
	return &PersonHandler{
		personStore: personStore,
	}
}

func (h *PersonHandler) HandlePutPerson(c *fiber.Ctx) error {
	var (
		personParams types.CreatePersonParams
		userID       = c.Params("id")
	)

	if err := c.BodyParser(&personParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	validationErrors := personParams.Validate()
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}
	person := types.UpdatedPersonFromParams(personParams)

	person, _ = h.personStore.Update(userID, &types.CreatePersonParams{
		Name:    person.Name,
		Age:     person.Age,
		Hobbies: person.Hobbies,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "one person updated",
		"person": person,
	})
}

func (h *PersonHandler) HandlePostPerson(c *fiber.Ctx) error {
	var personParams types.CreatePersonParams

	if err := c.BodyParser(&personParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	validationErrors := personParams.Validate()
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	person := types.NewPersonFromParams(personParams)
	h.personStore.Set(person)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "one person created",
		"person": person,
	})
}

func (h *PersonHandler) HandleGetPerson(c *fiber.Ctx) error {
	userID := c.Params("id")
	person, found := h.personStore.Get(userID)
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "one person found",
		"person": person,
	})
}

func (h *PersonHandler) HandleGetAllPersons(c *fiber.Ctx) error {
	persons := h.personStore.All()
	if len(*persons) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "no persons found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "all persons found",
		"persons": persons,
	})
}

func (h *PersonHandler) HandleDeletePerson(c *fiber.Ctx) error {
	userID := c.Params("id")
	_ = h.personStore.Delete(userID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "one person deleted",
	})
}
