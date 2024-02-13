package main

import (
	"net/http"

	"github.com/carlosgonzalez/learning-go/internal/handlers"
	"github.com/carlosgonzalez/learning-go/internal/middlewares"
	"github.com/carlosgonzalez/learning-go/internal/models"
	"github.com/carlosgonzalez/learning-go/internal/repositories"
	"github.com/carlosgonzalez/learning-go/pkg/validators"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate user model")
	}
	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		panic("failed to migrate post model")
	}

	//-------------------
	// Echo
	//-------------------

	e := echo.New()

	//Debug mode
	e.Debug = true

	//-------------------
	// Validator
	//-------------------
	e.Validator = validators.NewCustomValidator()

	//-------------------
	// Middlewares
	//-------------------

	// Stats
	s := middlewares.NewStats()
	e.Use(s.Process)

	//-------------------
	// Repositories
	//-------------------

	userRepository := repositories.NewUserRepository(db)

	//-------------------
	// Handlers
	//-------------------

	userHandler := handlers.NewUserHandler(userRepository)
	postHandler := handlers.NewPostHandler(db)

	//-------------------
	// Endpoints REST API
	//-------------------

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users", userHandler.GetAllUsers)
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUser)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	//-------------------
	// Endpoints RPC
	//-------------------

	e.POST("/fetch-posts", postHandler.CreatePost)

	e.Logger.Fatal(e.Start(":1323"))
}
