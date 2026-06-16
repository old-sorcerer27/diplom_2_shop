package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"strings"
	"time"
)

type ProcessedImage struct {
	Original  string `json:"original"`
	Medium    string `json:"medium"`
	Thumbnail string `json:"thumbnail"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
}

type ImageService struct {
	uploadDir   string
	productsDir string
	avatarsDir  string
}

func NewImageService() *ImageService {
	service := &ImageService{
		uploadDir:   AppConfig.UploadDir,
		productsDir: AppConfig.ProductsDir,
		avatarsDir:  AppConfig.AvatarsDir,
	}

	// Создаем необходимые директории
	service.ensureDirectories()

	return service
}

func (s *ImageService) ensureDirectories() {
	dirs := []string{
		s.uploadDir,
		s.productsDir,
		s.avatarsDir,
		filepath.Join(s.productsDir, "original"),
		filepath.Join(s.productsDir, "medium"),
		filepath.Join(s.productsDir, "thumbnail"),
		AppConfig.TempDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", dir, err)
		}
	}
}

func (s *ImageService) generateFilename(originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	name := uuid.New().String()
	date := time.Now().Format("2006/01/02")
	return filepath.Join(date, name+ext)
}

func (s *ImageService) isAllowedFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range AppConfig.AllowedFormats {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (s *ImageService) resizeImage(file io.Reader, width, height int) ([]byte, error) {
	// Декодируем изображение
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Получаем размеры
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Вычисляем новые размеры с сохранением пропорций
	if width > 0 && height > 0 {
		// Обрезаем до квадрата
		size := min(origWidth, origHeight)
		cropRect := image.Rect(
			(origWidth-size)/2,
			(origHeight-size)/2,
			(origWidth+size)/2,
			(origHeight+size)/2,
		)
		img = img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(cropRect)
		origWidth = size
		origHeight = size
	}

	// Масштабируем
	newWidth := width
	newHeight := height
	if newWidth == 0 {
		newWidth = origWidth
	}
	if newHeight == 0 {
		newHeight = origHeight
	}

	// Простое масштабирование (для продакшена используйте библиотеку типа "github.com/disintegration/imaging")
	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Кодируем результат
	buf := new(bytes.Buffer)
	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, result, &jpeg.Options{Quality: AppConfig.ImageQuality})
	case "png":
		err = png.Encode(buf, result)
	default:
		err = jpeg.Encode(buf, result, &jpeg.Options{Quality: AppConfig.ImageQuality})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	return buf.Bytes(), nil
}

func (s *ImageService) saveImage(data []byte, path string) error {
	fullPath := filepath.Join(s.productsDir, path)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// ProcessAndSave - основная функция для обработки и сохранения изображения
func (s *ImageService) ProcessAndSave(file io.Reader, originalName string, isAvatar bool) (*ProcessedImage, error) {
	// Проверяем формат
	if !s.isAllowedFormat(originalName) {
		return nil, fmt.Errorf("unsupported file format")
	}

	// Читаем файл
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if int64(len(data)) > AppConfig.ImageMaxSize {
		return nil, fmt.Errorf("file too large: max %d bytes", AppConfig.ImageMaxSize)
	}

	// Генерируем имя файла
	filename := s.generateFilename(originalName)
	_ = strings.TrimSuffix(filename, filepath.Ext(filename))

	var processed ProcessedImage
	processed.Filename = filename

	// Сохраняем оригинал
	originalPath := filepath.Join("original", filename)
	if err := s.saveImage(data, originalPath); err != nil {
		return nil, err
	}
	processed.Original = "/uploads/products/" + originalPath

	// Создаем среднее изображение
	mediumData, err := s.resizeImage(bytes.NewReader(data), AppConfig.MediumSize, 0)
	if err != nil {
		// Если не удалось создать среднее, используем оригинал
		processed.Medium = processed.Original
	} else {
		mediumPath := filepath.Join("medium", filename)
		if err := s.saveImage(mediumData, mediumPath); err != nil {
			processed.Medium = processed.Original
		} else {
			processed.Medium = "/uploads/products/" + mediumPath
		}
	}

	// Создаем миниатюру
	thumbnailData, err := s.resizeImage(bytes.NewReader(data), AppConfig.ThumbnailSize, AppConfig.ThumbnailSize)
	if err != nil {
		processed.Thumbnail = processed.Original
	} else {
		thumbnailPath := filepath.Join("thumbnail", filename)
		if err := s.saveImage(thumbnailData, thumbnailPath); err != nil {
			processed.Thumbnail = processed.Original
		} else {
			processed.Thumbnail = "/uploads/products/" + thumbnailPath
		}
	}

	processed.Size = int64(len(data))

	return &processed, nil
}

// DeleteImage - удаление изображения и всех его версий
func (s *ImageService) DeleteImage(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	// Извлекаем относительный путь
	relPath := strings.TrimPrefix(imageURL, "/uploads/products/")
	if relPath == imageURL {
		return nil
	}

	// Удаляем все версии
	paths := []string{
		filepath.Join(s.productsDir, "original", relPath),
		filepath.Join(s.productsDir, "medium", relPath),
		filepath.Join(s.productsDir, "thumbnail", relPath),
	}

	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Failed to delete file %s: %v\n", path, err)
		}
	}

	return nil
}

// DeleteGalleryImage - удаление изображения из галереи
func (s *ImageService) DeleteGalleryImage(imageURL string) error {
	return s.DeleteImage(imageURL)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
