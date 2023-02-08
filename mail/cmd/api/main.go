package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	PORT = "5004"
)

type Config struct {
	Mailer Mail
}

func main() {

	app := Config{
		Mailer: createMailer(),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	log.Println("[MAIL] Starting mail service at port", PORT)
	log.Fatal(server.ListenAndServe())
}

func createMailer() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	return Mail{
		Domain:      getFromEnv("MAIL_DOMAIN", ""),
		Host:        getFromEnv("MAIL_HOST", ""),
		Post:        port,
		Username:    getFromEnv("MAIL_USERNAME", ""),
		Password:    getFromEnv("MAIL_PASSWORD", ""),
		Encryption:  getFromEnv("MAIL_ENCRYPTION", "none"),
		FromName:    getFromEnv("MAIL_FROM_NAME", ""),
		FromAddress: getFromEnv("MAIL_FROM_ADDRESS", ""),
	}
}

func getFromEnv(env, def string) string {
	env_var := os.Getenv(env)
	if env_var == "" {
		env_var = def
	}
	return env_var
}
