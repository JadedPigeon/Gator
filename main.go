package main

import (
	"fmt"
	"log"

	"github.com/JadedPigeon/Gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config:", err)
	}
	fmt.Println("Current DB URL:", cfg.DBURL)
	fmt.Println("Current User:", cfg.CurrentUser)
	// Set a new user

	if err := cfg.SetUser("JadedPigeon"); err != nil {
		log.Fatal("error setting user:", err)
	}
	fmt.Println("Updated User:", cfg.CurrentUser)

	cfg, err = config.Read()
	if err != nil {
		log.Fatal("error reading config after update:", err)
	}
	fmt.Println("Config after update:", cfg)
}
