package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	sqc *ServerQueryClient
}

func (h *Handlers) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "i'm good xD",
	})
}

func (h *Handlers) intermission(c *gin.Context) {
	go h.sqc.exportMatchData()
	c.JSON(http.StatusOK, gin.H{
		"message": "good bro",
	})
}
