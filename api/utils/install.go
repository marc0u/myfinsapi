package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func InstallOnDocker(apiVersion string) {
	var cmdStr string
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
	-e "VIRTUAL_HOST=myfinsapi.aymconsulting.cl" \
	-e "VIRTUAL_PORT=%[1]v" \
	-e "LETSENCRYPT_HOST=myfinsapi.aymconsulting.cl" \
	-e "LETSENCRYPT_EMAIL=marco.urriola@gmail.com" \
	marc0u/go-launcher /app/myfinsapi`, strings.Replace(os.Getenv("API_PORT"), ":", "", 1))
	err := exec.Command("/bin/sh", "-c", cmdStr).Run()
	if err != nil {
		log.Fatalf("Docker couldn't be installed. Error code: %s", err)
	}
	fmt.Printf("> MyfinsAPI v%s successfully installed on Docker.\n", apiVersion)
}
