package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // sqlite database driver
	"github.com/marc0u/myfinsapi/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *fiber.App
}

func (server *Server) InitializeDB(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	switch Dbdriver {
	case "mysql":
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	case "postgres":
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	case "sqlite3":
		//DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DbName)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", Dbdriver)
		}
		server.DB.Exec("PRAGMA foreign_keys = ON")
	}
	if strings.ToLower(os.Getenv("DB_DEBUG")) == "true" {
		server.DB = server.DB.Debug()
	}
	if os.Getenv("DB_MIRROR") == "true" {
		err := server.MirrorProductionTables()
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		server.DB.AutoMigrate(&models.Transaction{}) //database migration
		server.DB.AutoMigrate(&models.Stock{})       //database migration

	}
}

func (server *Server) RunServer(addr, version string) {
	server.Router = fiber.New()
	server.Router.Use(cors.New())
	server.initializeRoutes(version)
	log.Fatal(server.Router.Listen(addr))
}
