package main

import (
	"fmt"

	"library_management/services"
)

func main() {
	lib := services.NewLibrary()
	_ = lib // wired to the console controller on Day 4
	fmt.Println("Library Management System - skeleton ready")
}
