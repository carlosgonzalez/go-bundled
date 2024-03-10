package main

import (
	"net/http"

	"github.com/carlosgonzalez/go-bundled/internal/handlers"
	"github.com/carlosgonzalez/go-bundled/internal/middlewares"
	"github.com/carlosgonzalez/go-bundled/internal/models"
	"github.com/carlosgonzalez/go-bundled/internal/repositories"
	"github.com/carlosgonzalez/go-bundled/pkg/validators"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Read from env variable
	dsn := "host=localhost user=postgres password=example dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
