package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
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
			fmt.Printf("  !  Cannot connect to %s database\n", Dbdriver)
			log.Fatal("  !  This is the error:", err)
		} else {
			fmt.Printf("  >  We are connected to the %s database\n", Dbdriver)
		}
	case "postgres":
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("  !  Cannot connect to %s database\n", Dbdriver)
			log.Fatal("  !  This is the error:", err)
		} else {
			fmt.Printf("  >  We are connected to the %s database\n", Dbdriver)
		}
	case "sqlite3":
		//DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DbName)
		if err != nil {
			fmt.Printf("  !  Cannot connect to %s database\n", Dbdriver)
			log.Fatal("  !  This is the error:", err)
		} else {
			fmt.Printf("  >  We are connected to the %s database\n", Dbdriver)
		}
		server.DB.Exec("PRAGMA foreign_keys = ON")
	}
	if strings.ToLower(os.Getenv("DB_DEBUG")) == "true" {
		server.DB = server.DB.Debug()
	}
	server.DB.AutoMigrate(&models.Transaction{}) //database migration
	server.DB.AutoMigrate(&models.Stock{})       //database migration
	if len(os.Args) > 1 && os.Args[1] == "mirror" {
		fmt.Println("  >  Mirroring online tables...")
		err := server.MirrorProductionTables()
		if err != nil {
			fmt.Println("  !  Error Mirroring databases.")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("  >  Tables Mirrored successfuly.")
	}
}

func (s *Server) RunServer(addr, version string) {
	s.Router = fiber.New()
	s.initializeMiddlewares()
	s.initializeRoutes(version)
	log.Fatal(s.Router.Listen(addr))
}
