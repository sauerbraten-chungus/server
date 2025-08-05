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
		log.Printf("Error loading .env file: %s\n", err)
	}

	playerServiceIP := os.Getenv("PLAYER_SERVICE_IP")
	gameServerIP := os.Getenv("GAME_SERVER_IP")
	authServiceIP := os.Getenv("AUTH_SERVICE_IP")
	gameServerPort, err := strconv.Atoi(os.Getenv("GAME_SERVER_PORT"))
	secretKey := os.Getenv("SECRET_CHUNGUS")
	apiKey := os.Getenv("CHUNGUS_KEY")
	if err != nil {
		log.Fatalf("Error loading game server port: %s\n", err)
	}

	sqc, err := NewServerQueryClient(gameServerIP, playerServiceIP, authServiceIP, apiKey, gameServerPort)
	if err != nil {
		log.Fatalf("Error starting server query client: %s\n", err)
	}

	handlers := &Handlers{
		sqc: sqc,
	}

	web := gin.Default()
	web.GET("/health", handlers.health)
	protected := web.Group("/")
	protected.Use(JWTAuthMiddleware(secretKey))
	{
		protected.GET("intermission", handlers.intermission)
	}

	err = web.Run()
	if err != nil {
		log.Fatalf("Error starting server lol\n")
	}
}
