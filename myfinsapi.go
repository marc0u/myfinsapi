package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/marc0u/myfinsapi/api/handlers"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func main() {
	apiVersion := "2.3.0"
	var err error
	var cmdStr string
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			fmt.Println("> Removing previous version...")
			exec.Command("/bin/sh", "-c", "docker container rm myfinsapi").Run()
			fmt.Printf("> Installing MyfinsAPI v%s on Docker.\n", apiVersion)
			cmdStr = fmt.Sprintf(`docker run -d \
		--name myfinsapi \
		--network myfins \
		--restart=unless-stopped \
		-v $PWD:/app \
		-e TZ=America/Santiago \
		-p %[1]v:%[1]v \
		marc0u/go-launcher /app/myfinsapi`, strings.Replace(os.Getenv("API_PORT"), ":", "", 1))
			err = exec.Command("/bin/sh", "-c", cmdStr).Run()
			if err != nil {
				log.Fatalf("Docker couldn't be installed. Error code: %s", err)
			}
			fmt.Printf("> MyfinsAPI v%s successfully installed on Docker.\n", apiVersion)
		}
	} else {
		fmt.Printf("---------- MyfinsAPI v%s ----------\n", apiVersion)
		var server = handlers.Server{}
		// Initialize DB
		server.InitializeDB(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
		// Initialize Server
		server.RunServer(os.Getenv("API_PORT"), apiVersion)
	}
}
