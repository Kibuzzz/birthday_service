package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Server) AllEmployees(c *fiber.Ctx) error {
	users, err := h.UserStorage.GetAll()
	if err != nil {
		return c.SendString("db error")
	}
	return c.JSON(users)
}

func (h *Server) Employee(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.SendString("error email")
	}
	user, err := h.UserStorage.GetByEmail(email)
	if err != nil {
		return c.SendString("user not found")
	}
	return c.JSON(user)
}

func (h *Server) DeleteEmployee(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.SendString("error email")
	}
	err := h.UserStorage.DeleteEmp(email)
	if err != nil {
		return c.JSON("failed to delete user")
	}
	return c.SendStatus(http.StatusOK)
}
