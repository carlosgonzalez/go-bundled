package main

import (
	"net/http"

	"github.com/carlosgonzalez/learning-go/handlers"
	"github.com/carlosgonzalez/learning-go/middlewares"
	"github.com/carlosgonzalez/learning-go/models"
	"github.com/carlosgonzalez/learning-go/services"

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
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Post{})

	//-------------------
	// Echo
	//-------------------

	e := echo.New()

	//Debug mode
	e.Debug = true

	//-------------------
	// Validator
	//-------------------
	e.Validator = services.NewCustomValidator()

	//-------------------
	// Middlewares
	//-------------------

	// Stats
	s := middlewares.NewStats()
	e.Use(s.Process)

	//-------------------
	// Handlers
	//-------------------

	userHandler := handlers.NewUserHandler(db)
	postHandler := handlers.NewPostHandler(db)

	//-------------------
	// Endpoints
	//-------------------

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users", userHandler.GetAllUsers)
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUser)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	e.POST("/fetch-posts", postHandler.CreatePost)

	e.Logger.Fatal(e.Start(":1323"))
}
