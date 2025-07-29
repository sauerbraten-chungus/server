package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file\n")
	}

	playerServiceIP := os.Getenv("PLAYER_SERVICE_IP")
	gameServerIP := os.Getenv("GAME_SERVER_IP")
	gameServerPort, err := strconv.Atoi(os.Getenv("GAME_SERVER_PORT"))
	if err != nil {
		log.Fatalf("Error loading game server port\n")
	}

	sqc, err := NewServerQueryClient(gameServerIP, playerServiceIP, gameServerPort)
	if err != nil {
		log.Fatalf("Error starting server query client")
	}

	handlers := &Handlers{
		sqc: sqc,
	}

	web := gin.Default()
	web.GET("/health", handlers.health)
	web.GET("/intermission", handlers.intermission)

	err = web.Run()
	if err != nil {
		log.Fatalf("Error starting server lol\n")
	}
}
