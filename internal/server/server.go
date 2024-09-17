package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"go-api/internal/dal"
	"go-api/internal/database"
	"go-api/internal/services"
)

type FiberServer struct {
	DB *database.Service
	*fiber.App
	User  services.UserService
	Task  services.TaskService
	DAO   *dal.Query
	Store *session.Store
}

var server *FiberServer

func New() *FiberServer {
	if server != nil {
		return server
	}

	db := database.New()
	gorm := db.GetORM()
	dao := dal.Use(gorm)
	storage := sqlite3.New()
	store := session.New(session.Config{
		Storage: storage,
	})

	userService := services.NewUserService(dao)
	taskService := services.NewTaskService(dao)

	server = &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "go-api",
			AppName:      "go-api",
		}),
		DB:    db,
		User:  userService,
		Task:  taskService,
		DAO:   dao,
		Store: store,
	}

	return server
}
