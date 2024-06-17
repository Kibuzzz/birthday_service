package server

import (
	model "birtday_service/models/user"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) AllEmployees(c *fiber.Ctx) error {

	s.Logger.Info("all users handler")

	users, err := s.UserStorage.GetAll()
	noRows := errors.Is(err, model.ErrorNotFound)
	switch {
	case err != nil && !noRows:
		s.Logger.Debugw("internal error", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal error"})
	case noRows:
		s.Logger.Debug("no users")
		return c.JSON(fiber.Map{"message": "no users"})
	}
	return c.JSON(users)
}

func (s *Server) DeleteEmployee(c *fiber.Ctx) error {

	s.Logger.Info("delete emp handler")

	type DeleteRequest struct {
		Email string `json:"email"`
	}

	var delete DeleteRequest
	if err := c.BodyParser(&delete); err != nil {
		s.Logger.Debugw("bad request", "request", delete, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "bad request"})
	}

	if delete.Email == "" {
		s.Logger.Debugw("bad request", "error", "empty email")
		return c.Status(fiber.StatusBadRequest).SendString("error email")
	}

	err := s.UserStorage.DeleteEmp(delete.Email)
	if err != nil {
		s.Logger.Debugw("failed to delte user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to delete user"})
	}
	return c.SendStatus(http.StatusOK)
}
