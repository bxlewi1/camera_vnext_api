package main

import (
	"camera_vnext_api/config"
	"camera_vnext_api/routes"
	"log"
)

func main() {
	log.Fatalln(routes.App.Listen(config.Instance.Port))
}
