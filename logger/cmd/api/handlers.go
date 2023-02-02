package main

import (
	"log"
	"logger/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

type requestPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) writeLog(c *gin.Context) {

	var req requestPayload

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	entry := data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := app.Models.LogEntry.Insert(entry)

	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
