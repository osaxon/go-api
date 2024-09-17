package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go-api/internal/models"
	"os"
)

func (s *FiberServer) RegisterFiberRoutes() {

	auth := s.App.Group("auth")
	api := s.App.Group("api")
	api.Use(AuthMiddleware)
	api.Get("/", s.HelloWorldHandler)
	api.Get("/health", s.healthHandler)
	api.Post("/users", s.createUser)
	api.Get("/users/:id", s.getUser)
	api.Post("/tasks", s.createTask)
	api.Get("/tasks/:id", s.getTask)
	api.Get("/tasks", s.getTasks)

	auth.Post("/register", s.register)
	auth.Post("/login", s.login)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}
	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.DB.Health())
}

func (s *FiberServer) createUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	resp, _ := s.User.CreateUser(user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": resp,
	})
}

func (s *FiberServer) getUser(c *fiber.Ctx) error {
	q := c.QueryBool("expand")
	id := c.Params("id")
	sess, err := s.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println(sess)
	resp, _ := s.User.GetUserById(id, q)
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (s *FiberServer) createTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	resp, _ := s.Task.CreateTask(&task)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"task": resp,
	})
}

func (s *FiberServer) getTask(c *fiber.Ctx) error {
	id := c.Params("id")

	resp, _ := s.Task.GetTaskById(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"task": resp,
	})
}

func (s *FiberServer) getTasks(c *fiber.Ctx) error {
	tasks, err := s.Task.GetTasks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasks": tasks,
	})
}

func (s *FiberServer) register(c *fiber.Ctx) error {

	var registerRequest struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := s.User.RegisterUser(registerRequest.Name, registerRequest.Email, registerRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"user": user,
	})
}

func (s *FiberServer) login(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := s.User.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"token": token,
	})
}

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the token signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("CRYPTO_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid Token")
	}

	// Extract user information from the token and store it in the context
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		user := getUserFromToken(claims)
		c.Locals("user", user)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).SendString("Invalid Token Claims")
}

func getUserFromToken(claims jwt.MapClaims) models.User {
	user := models.User{
		BaseModel: models.BaseModel{
			ID: claims["sub"].(string),
		},
		Email: claims["email"].(string),
	}
	return user
}
