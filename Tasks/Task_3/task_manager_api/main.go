package main

import (
	"log"
	"task_manager_api/router"
)

func main() {
	r := router.SetupRouter()
	log.Fatal(r.Run(":8080"))
}
