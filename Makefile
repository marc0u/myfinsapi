hello:
	echo "Hello"

build:
	env GOOS=linux GOARCH=amd64 GOARM=7 go build -o bin/myfinsapi -v myfinsapi.go

run:
	go run myfinsapi.go