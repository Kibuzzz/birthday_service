package server

import (
	subscription "birtday_service/models/subscription"
	usr "birtday_service/models/user"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	subDateFormat = "15:04"
)

func (s *Server) Subscribe(c *fiber.Ctx) error {
	s.Logger.Info("subscribe handler")

	subsID, ok := c.Locals("id").(int)
	if !ok {
		s.Logger.Error("id not found")
		return c.JSON(fiber.Map{"message": "internal error", "error": subsID})
	}

	type SubRequest struct {
		EmployeeID        int    `json:"id"`
		TimeBeforBirthday string `json:"time"` // время до отправления уведомления
	}

	var subRequest SubRequest
	if err := c.BodyParser(&subRequest); err != nil {
		s.Logger.Errorw("bad requset", "body", subRequest)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := s.UserStorage.GetByID(subRequest.EmployeeID)
	userNotFound := errors.Is(err, usr.ErrorNotFound)
	switch {
	case err != nil && !userNotFound:
		s.Logger.Errorw("inernal error", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	case userNotFound:
		s.Logger.Info("user not found")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "employee not found"})
	}

	timeBeforeBirthday, err := time.Parse(subDateFormat, subRequest.TimeBeforBirthday)
	if err != nil {
		s.Logger.Info("bad time formatting")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "bad time. wanted hh:mm format", "error": err})
	}

	hours := time.Duration(timeBeforeBirthday.Hour())
	minutes := time.Duration(timeBeforeBirthday.Minute())

	thisYearBirthday := time.Date(time.Now().Year(), user.Birthday.Month(), user.Birthday.Day(), user.Birthday.Hour(), user.Birthday.Minute(), 0, 0, subscription.Location)
	notificationTime := thisYearBirthday.Add(-time.Hour * hours)
	notificationTime = notificationTime.Add(-time.Minute * minutes)

	_, err = s.SubsStorage.GetSub(subsID, subRequest.EmployeeID)
	noRows := errors.Is(err, subscription.ErrorNoRows)
	if err != nil && !noRows {
		s.Logger.Errorw("internal error", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal error"})
	}

	if noRows {
		err = s.SubsStorage.Sub(subsID, subRequest.EmployeeID, notificationTime)
		if err != nil {
			s.Logger.Errorw("subscribing error", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal error"})
		}
	} else {
		err = s.SubsStorage.UpdateSub(subsID, subRequest.EmployeeID, notificationTime)
		if err != nil {
			s.Logger.Errorw("updating subscription error", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal error"})
		}
	}
	return c.JSON(fiber.Map{"notification time": timeBeforeBirthday.String()})
}

func (s *Server) UnSubscribe(c *fiber.Ctx) error {

	s.Logger.Info("unsubscribe handler")

	subsID, ok := c.Locals("id").(int)
	if !ok {
		s.Logger.Error("id not found")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	type UnSubRequest struct {
		EmployeeID int `json:"id"`
	}

	var unSub UnSubRequest
	if err := c.BodyParser(&unSub); err != nil {
		s.Logger.Debugw("bad request", "request", unSub)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := s.SubsStorage.UnSub(subsID, unSub.EmployeeID)
	subNotFound := errors.Is(err, subscription.ErrorNoRows)

	switch {
	case err != nil && !subNotFound:
		s.Logger.Errorw("internal error", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal error"})
	case subNotFound:
		s.Logger.Debugf("no subscription with subsID: %d and celebratorID: %d", subsID, unSub.EmployeeID)
		return c.JSON(fiber.Map{"message": "no such subcription"})
	}
	return c.JSON(fiber.Map{"message": "success"})
}

func (s *Server) SubsByID(c *fiber.Ctx) error {

	s.Logger.Info("all user subscriptions handler")

	id, ok := c.Locals("id").(int)
	if !ok {
		s.Logger.Errorw("auth error", "error", "failed to get id from c.Locals")
		return c.JSON(fiber.Map{"message": "internal error", "error": id})
	}

	subs, err := s.SubsStorage.GetSubsByID(id)
	isEmpty := errors.Is(err, subscription.ErrorNoRows)

	switch {
	case err != nil && !isEmpty:
		s.Logger.Errorw("internal error", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	case isEmpty:
		s.Logger.Debug("no subscriptions")
		return c.JSON(fiber.Map{"message": "no subscriptions"})
	}
	return c.JSON(subs)
}
