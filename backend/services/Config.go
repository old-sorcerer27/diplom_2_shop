package services

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	DBPath string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	OwnerEmail   string

	JWTSecret string
}

var AppConfig *Config

func LoadConfig() error {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	AppConfig = &Config{
		Port:         os.Getenv("BACKEND_PORT"),
		DBPath:       os.Getenv("DB_PATH"),
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     smtpPort,
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
		OwnerEmail:   os.Getenv("OWNER_EMAIL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}

	return nil
}
