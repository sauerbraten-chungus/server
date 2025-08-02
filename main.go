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
	authServiceIP := os.Getenv("AUTH_SERVICE_IP")
	gameServerPort, err := strconv.Atoi(os.Getenv("GAME_SERVER_PORT"))
	secretKey := os.Getenv("SECRET_CHUNGUS")
	if err != nil {
		log.Fatalf("Error loading game server port\n")
	}

	sqc, err := NewServerQueryClient(gameServerIP, playerServiceIP, authServiceIP, gameServerPort)
	if err != nil {
		log.Fatalf("Error starting server query client")
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
