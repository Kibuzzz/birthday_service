package server

import (
	usr "birtday_service/models/user"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	birthdayFormat = "02.01.2006"
)

func (s *Server) Register(c *fiber.Ctx) error {

	s.Logger.Info("registrer handler")

	type RegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday string `json:"birthday"`
	}

	var register RegisterRequest
	err := c.BodyParser(&register)
	if err != nil {
		s.Logger.Infow("bad request", "request", register, "error", err)
		return fmt.Errorf("bad user: %w", err)
	}

	_, err = s.UserStorage.GetByEmail(register.Email)
	userNotFound := errors.Is(err, usr.ErrorNotFound)
	switch {
	case err != nil && !userNotFound:
		s.Logger.Errorw("internal error", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	case !userNotFound:
		s.Logger.Infow("duplicat user", "user email", register.Email)
		return c.JSON(fmt.Sprintf("user with email %s exists", register.Email))
	}

	birthdayDate, err := time.Parse(birthdayFormat, register.Birthday)
	if err != nil {
		s.Logger.Infow("bad time formatting", "time", register.Birthday)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "bad birthday. wanted day.month.year"})
	}

	err = s.UserStorage.AddEmp(register.Name, register.Email, register.Password, birthdayDate)
	if err != nil {
		s.Logger.Errorw("failed to add employee", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to registrate"})
	}
	return c.SendStatus(http.StatusOK)
}

var jwtSecret = []byte("jwtSecret")

func (s *Server) Login(c *fiber.Ctx) error {

	s.Logger.Info("login handler")

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var login LoginRequest
	if err := c.BodyParser(&login); err != nil {
		s.Logger.Infow("bad request", "request", login)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	email := login.Email
	pass := login.Password
	user, err := s.UserStorage.GetByEmail(email)
	userNotFound := errors.Is(err, usr.ErrorNotFound)
	switch {
	case err != nil && !userNotFound:
		s.Logger.Errorw("internal error", "request", login, "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	case userNotFound:
		s.Logger.Infow("no such user", "request", login)
		return c.JSON(fiber.Map{"message": "user not found"})
	}

	if pass != user.Password {
		s.Logger.Debugw("wrong password", "password", pass)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "wrong password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		s.Logger.Errorw("failed to sign token", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	tokenCookie := new(fiber.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Value = signedToken
	tokenCookie.Expires = time.Now().Add(72 * time.Hour)
	tokenCookie.HTTPOnly = true
	c.Cookie(tokenCookie)

	return c.JSON(fiber.Map{"message": "Success login", "token": signedToken})
}

func (s *Server) Logout(c *fiber.Ctx) error {

	s.Logger.Info("logout handler")

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Successfully logged out"})
}

func (s *Server) AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Missing or invalid JWT"})
	}
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "JWT has expired"})
			}
			idFloat, ok := claims["id"].(float64)
			if !ok {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "bad id"})
			}
			id := int(idFloat)
			c.Locals("email", claims["email"])
			c.Locals("id", id)
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT claims"})
	}
	return c.Next()
}
