package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/like2foxes/nirlir/internal/router"
)

func main() {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	router.Start(jwtSecret)
}
