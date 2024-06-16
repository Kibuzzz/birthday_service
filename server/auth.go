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

func (h *Server) Register(c *fiber.Ctx) error {

	type RegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday string `json:"birthday"`
	}

	var register RegisterRequest
	err := c.BodyParser(&register)
	if err != nil {
		return fmt.Errorf("bad user: %w", err)
	}

	user, err := h.UserStorage.GetByEmail(register.Email)
	userNotFound := errors.Is(err, usr.ErrorNotFound)
	switch {
	case err != nil && !userNotFound:
		return c.SendStatus(fiber.StatusInternalServerError)
	case user != usr.User{}:
		return c.JSON(fmt.Sprintf("user with email %s exists", register.Email))
	}

	birthdayDate, err := time.Parse(birthdayFormat, register.Birthday)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "bad birthday. wanted day.month.year"})
	}

	err = h.UserStorage.AddEmp(register.Name, register.Email, register.Password, birthdayDate)
	if err != nil {
		return c.JSON(fmt.Sprintf("faile to add user: %v", err))
	}
	return c.SendStatus(http.StatusOK)
}

var jwtSecret = []byte("jwtSecret")

func (h *Server) Login(c *fiber.Ctx) error {

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var login LoginRequest
	if err := c.BodyParser(&login); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	email := login.Email
	pass := login.Password
	user, err := h.UserStorage.GetByEmail(email)
	userNotFound := errors.Is(err, usr.ErrorNotFound)
	switch {
	case err != nil && !userNotFound:
		return c.SendStatus(fiber.StatusInternalServerError)
	case userNotFound:
		return c.JSON(fiber.Map{"message": "user not found"})
	}

	if pass != user.Password {
		fmt.Println(login)
		fmt.Println(email, pass)
		fmt.Println(user.Email, user.Password)
		return c.JSON(fiber.Map{"message": "wrong password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
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

func (h *Server) Logout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Successfully logged out"})
}

func (h *Server) AuthMiddleware(c *fiber.Ctx) error {
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
