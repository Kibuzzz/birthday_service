package server

import (
	sub "birtday_service/models/subscription"
	usr "birtday_service/models/user"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	subDateFormat = "15:04"
)

func (h *Server) Subscribe(c *fiber.Ctx) error {
	subsID, ok := c.Locals("id").(int)
	if !ok {
		return c.JSON(fiber.Map{"message": "internal error", "error": subsID})
	}

	type SubRequest struct {
		EmployeeID        int    `json:"id"`
		TimeBeforBirthday string `json:"time"` // время до отправления уведомления
	}

	var subRequest SubRequest
	if err := c.BodyParser(&subRequest); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := h.UserStorage.GetByID(subRequest.EmployeeID)
	userNotFound := errors.Is(err, usr.ErrorNotFound)
	switch {
	case err != nil && !userNotFound:
		return c.SendStatus(fiber.StatusInternalServerError)
	case userNotFound:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "employee not found"})
	}

	timeBeforeBirthday, err := time.Parse(subDateFormat, subRequest.TimeBeforBirthday)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "bad time. wanted hh:mm format", "error": err})
	}

	err = h.SubsStorage.Sub(subsID, subRequest.EmployeeID, user.Birthday, timeBeforeBirthday)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"time": timeBeforeBirthday.String()})
}

func (h *Server) UnSubscribe(c *fiber.Ctx) error {
	subsID, ok := c.Locals("id").(int)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	type UnSubRequest struct {
		EmployeeID int `json:"id"`
	}

	var unSubRequest UnSubRequest
	if err := c.BodyParser(&unSubRequest); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := h.SubsStorage.UnSub(subsID, unSubRequest.EmployeeID)
	subNotFound := errors.Is(err, sub.ErrorNotFound)
	switch {
	case err != nil && !subNotFound:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal error"})
	case subNotFound:
		return c.JSON(fiber.Map{"message": "no such subcription"})
	}
	return c.JSON(fiber.Map{"message": "success"})
}

func (h *Server) SubsByID(c *fiber.Ctx) error {
	id, ok := c.Locals("id").(int)
	if !ok {
		return c.JSON(fiber.Map{"message": "internal error", "error": id})
	}
	subs, err := h.SubsStorage.GetSubsByID(id)
	isEmpty := errors.Is(err, sub.ErrorEmpty)
	switch {
	case err != nil && !isEmpty:
		return c.SendStatus(fiber.StatusInternalServerError)
	case isEmpty:
		return c.JSON(fiber.Map{"message": "no subscriptions"})
	}
	return c.JSON(subs)
}
