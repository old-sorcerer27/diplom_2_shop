package services

import (
	"log"
	"os"
	"path/filepath"
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

	ImageMaxSize   int64 // bytes
	ImageQuality   int   // 1-100
	ThumbnailSize  int   // pixels
	MediumSize     int   // pixels
	LargeSize      int   // pixels
	AllowedFormats []string

	UploadDir   string
	ProductsDir string
	AvatarsDir  string
	TempDir     string
}

var AppConfig *Config

func LoadConfig() error {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	imageMaxSize, _ := strconv.ParseInt(os.Getenv("IMAGE_MAX_SIZE"), 10, 64)
	imageQuality, _ := strconv.Atoi(os.Getenv("IMAGE_QUALITY"))
	thumbnailSize, _ := strconv.Atoi(os.Getenv("THUMBNAIL_SIZE"))
	mediumSize, _ := strconv.Atoi(os.Getenv("MEDIUM_SIZE"))
	largeSize, _ := strconv.Atoi(os.Getenv("LARGE_SIZE"))

	AppConfig = &Config{
		Port:           os.Getenv("BACKEND_PORT"),
		DBPath:         os.Getenv("DB_PATH"),
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       smtpPort,
		SMTPUser:       os.Getenv("SMTP_USER"),
		SMTPPassword:   os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:       os.Getenv("SMTP_FROM"),
		OwnerEmail:     os.Getenv("OWNER_EMAIL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		ImageMaxSize:   imageMaxSize,
		ImageQuality:   imageQuality,
		ThumbnailSize:  thumbnailSize,
		MediumSize:     mediumSize,
		LargeSize:      largeSize,
		AllowedFormats: []string{".jpg", ".jpeg", ".png", ".webp", ".gif"},
		UploadDir:      os.Getenv("UPLOAD_DIR"),
		ProductsDir:    os.Getenv("PRODUCTS_DIR"),
		AvatarsDir:     os.Getenv("AVATARS_DIR"),
		TempDir:        os.Getenv("TEMP_DIR"),
	}

	log.Printf("Loaded config: %+v\n", AppConfig)

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func SetupUploadDirectories() {
	dirs := []string{
		AppConfig.UploadDir,
		AppConfig.ProductsDir,
		AppConfig.AvatarsDir,
		filepath.Join(AppConfig.ProductsDir, "original"),
		filepath.Join(AppConfig.ProductsDir, "medium"),
		filepath.Join(AppConfig.ProductsDir, "thumbnail"),
		AppConfig.TempDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Failed to create directory %s: %v", dir, err)
		}
	}

	log.Println("Upload directories initialized")
}
