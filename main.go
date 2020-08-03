package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/jinzhu/gorm"

	db "github.com/jhankes/sample/db"
	group "github.com/jhankes/sample/model"
	user "github.com/jhankes/sample/model"

	"github.com/jhankes/sample/model"
)

// Create the fiber routes for the API
func setupRoutes(app *fiber.App) {

	app.Use(middleware.Logger())

	// users routes
	app.Get("/api/v1/users/:userid", user.GetUser)
	app.Post("/api/v1/users", user.NewUser)
	app.Put("/api/v1/users/:userid", user.UpdateUser)
	app.Delete("/api/v1/users/:userid", user.DeleteUser)

	// groups routes
	app.Get("/api/v1/groups/:name", group.GetGroup)
	app.Post("/api/v1/groups", group.NewGroup)
	app.Put("/api/v1/groups/:name", group.UpdateGroup)
	app.Delete("/api/v1/groups/:name", group.DeleteGroup)
}

// Configure the postgres database with the default settings or use the env vars if available
func initDatabase() {
	dbConfig := DatabaseConfig{Host: "localhost", Port: "5432", Username: "sample", DbName: "sample", Password: "sample", SslMode: "disable"}
	if os.Getenv("SAMPLE_HOST") != "" {
		dbConfig.Host = os.Getenv("SAMPLE_HOST")
	}
	if os.Getenv("SAMPLE_PORT") != "" {
		dbConfig.Host = os.Getenv("SAMPLE_PORT")
	}
	if os.Getenv("SAMPLE_DBNAME") != "" {
		dbConfig.Host = os.Getenv("SAMPLE_DBNAME")
	}
	if os.Getenv("SAMPLE_USER") != "" {
		dbConfig.Host = os.Getenv("SAMPLE_USER")
	}
	if os.Getenv("SAMPLE_PASS") != "" {
		dbConfig.Host = os.Getenv("SAMPLE_PASS")
	}
	if os.Getenv("SAMPLE_SSLMODE") != "" {
		dbConfig.Host = os.Getenv("SAMPLE_SSLMODE")
	}
	var err error
	cxn := "host=" + dbConfig.Host + " port=" + dbConfig.Port + " user=" + dbConfig.Username +
		" dbname=" + dbConfig.DbName + " password=" + dbConfig.Password + " sslmode=" + dbConfig.SslMode
	db.DBConn, err = gorm.Open("postgres", cxn)
	if err != nil {
		panic(err)
	}
	log.Println("Database connection opened successfully...")
	db.DBConn.AutoMigrate(&model.User{}, &model.Group{}, &model.Association{})
	log.Println("Database migrated User, Group, Association")
}

// Main
func main() {
	app := fiber.New()

	initDatabase()

	setupRoutes(app)
	app.Listen(3000)

	defer db.DBConn.Close()
}

// DatabaseConfig : Minimal postgres database config
type DatabaseConfig struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
	SslMode  string
}
