package main

import (
	"log"
	"os"

	"./webServer"
)

func getPort () string {
	if (os.Getenv("PORT") != "") {
		return os.Getenv("PORT")
	}

	defaultPort := "8080"

	log.Println("Missing HTTP Server Listening port using", defaultPort, "as default port")

	return defaultPort
}

func main() {
	log.Fatal(webServer.Run(getPort()))
}


