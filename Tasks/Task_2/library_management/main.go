package main

import (
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
)

func main() {
	lib := services.NewLibrary()

	// Seed members (the interface has no AddMember).
	lib.AddMember(models.Member{ID: 1, Name: "Abel"})
	lib.AddMember(models.Member{ID: 2, Name: "Sara"})

	controllers.NewLibraryController(lib).Run()
}
