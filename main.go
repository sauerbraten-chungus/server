package main

import (
	"log"
	"os"
	"strconv"

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
	sqc.exportMatchData()
}
