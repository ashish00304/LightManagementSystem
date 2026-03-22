package main

import (
	"fmt"
	"light-management/routes"
)

func main() {
	r := routes.SetupRouter()
	fmt.Println("Starting Light Management Server on :8080...")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Startup failed: %v\n", err)
	}
}
