package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) SendMail(c *gin.Context) {

	var requestPayload struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	if err := c.BindJSON(&requestPayload); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	if err := app.Mailer.SendSMTPMessage(msg); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Mail sent successfully to " + requestPayload.To,
	})
}
