package database

import (
	"log"
	"os"
	"store/database/source"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func ConnectDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("./store.db"), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных: " + err.Error())
	}
}

func MigrateDB() {
	err := DB.AutoMigrate(&source.User{}, &source.Product{}, &source.Order{}, &source.OrderItems{})
	if err != nil {
		log.Printf("Ошибка миграции: %v", err)
		return
	}
}

func InsertSampleData() {
	result := DB.Create(&source.User{
		Username: "zhizha",
		Password: "password",
	})
	if result.Error != nil {
		log.Printf("Ошибка при создании тестового пользователя: %v", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("Запись не была создана, возможно из-за уникального ограничения на имя пользователя")
		return
	}

	products := []source.Product{
		{
			Name:        "Ноутбук",
			Description: "Мощный ноутбук для работы и игр с процессором Intel Core i7, 16GB RAM, 512GB SSD",
			Price:       59999,
			Stock:       10,
			ImageURL:    "/uploads/products/original/laptop.jpg",
			Category:    "Ноутбуки",
		},
		{
			Name:        "Смартфон",
			Description: "Современный смартфон с отличной камерой 108MP, AMOLED дисплеем 120Hz",
			Price:       29999,
			Stock:       20,
			ImageURL:    "/uploads/products/original/smartphone.jpg",
			Category:    "Смартфоны",
		},
		{
			Name:        "Наушники",
			Description: "Беспроводные наушники с активным шумоподавлением и временем работы до 30 часов",
			Price:       4999,
			Stock:       30,
			ImageURL:    "/uploads/products/original/headphones.jpg",
			Category:    "Аудио",
		},
		{
			Name:        "Клавиатура",
			Description: "Механическая клавиатура с RGB подсветкой и красными switches",
			Price:       3999,
			Stock:       15,
			ImageURL:    "/uploads/products/original/keyboard.jpg",
			Category:    "Периферия",
		},

		// Новые товары
		{
			Name:        "Игровая мышь Logitech G502",
			Description: "Высокоточная игровая мышь с 11 программируемыми кнопками и сенсором 25K DPI",
			Price:       5490,
			Stock:       25,
			ImageURL:    "/uploads/products/original/mouse.jpg",
			Category:    "Периферия",
		},
		{
			Name:        "Монитор 27\" 4K",
			Description: "27-дюймовый 4K монитор с IPS матрицей, 100% sRGB и HDR400",
			Price:       32990,
			Stock:       8,
			ImageURL:    "/uploads/products/original/monitor.jpg",
			Category:    "Мониторы",
		},
		{
			Name:        "SSD накопитель 1TB",
			Description: "Скоростной SSD NVMe M.2 с чтением до 7000 MB/s",
			Price:       7990,
			Stock:       50,
			ImageURL:    "/uploads/products/original/ssd.jpg",
			Category:    "Комплектующие",
		},
		{
			Name:        "Видеокарта RTX 4060",
			Description: "Видеокарта для игр и работы с трассировкой лучей, 8GB GDDR6",
			Price:       35990,
			Stock:       5,
			ImageURL:    "/uploads/products/original/gpu.jpg",
			Category:    "Комплектующие",
		},
		{
			Name:        "Процессор Intel i9-13900K",
			Description: "Топовый процессор для максимальной производительности, 24 ядра",
			Price:       58990,
			Stock:       3,
			ImageURL:    "/uploads/products/original/cpu.jpg",
			Category:    "Комплектующие",
		},
		{
			Name:        "Оперативная память 32GB DDR5",
			Description: "Двухканальный комплект RAM с частотой 6000MHz",
			Price:       11990,
			Stock:       20,
			ImageURL:    "/uploads/products/original/ram.jpg",
			Category:    "Комплектующие",
		},
		{
			Name:        "Блок питания 850W Gold",
			Description: "Блок питания с сертификатом 80+ Gold и модульными кабелями",
			Price:       8990,
			Stock:       12,
			ImageURL:    "/uploads/products/original/psu.jpg",
			Category:    "Комплектующие",
		},
		{
			Name:        "Игровой стул",
			Description: "Эргономичное кресло с поддержкой поясницы и регулировкой подлокотников",
			Price:       15990,
			Stock:       7,
			ImageURL:    "/uploads/products/original/chair.jpg",
			Category:    "Мебель",
		},
		{
			Name:        "Веб-камера 4K",
			Description: "Веб-камера с 4K разрешением, автофокусом и встроенным микрофоном",
			Price:       8990,
			Stock:       15,
			ImageURL:    "/uploads/products/original/webcam.jpg",
			Category:    "Периферия",
		},
		{
			Name:        "Микрофон Yeti",
			Description: "Студийный USB микрофон с кардиоидной диаграммой направленности",
			Price:       12990,
			Stock:       10,
			ImageURL:    "/uploads/products/original/microphone.jpg",
			Category:    "Аудио",
		},
		{
			Name:        "Колонки 2.1",
			Description: "Акустическая система с сабвуфером общей мощностью 80W",
			Price:       6990,
			Stock:       18,
			ImageURL:    "/uploads/products/original/speakers.jpg",
			Category:    "Аудио",
		},
	}
	for _, p := range products {
		result := DB.Create(&p)
		if result.Error != nil {
			log.Printf("Ошибка при создании тестового продукта: %v", result.Error)
			return
		}

		if result.RowsAffected == 0 {
			log.Printf("Запись не была создана, возможно из-за уникального ограничения на имя продукта: %s", p.Name)
			return
		}
		log.Printf("Тестовый продукт создан: %v", p.Name)
	}

}

func createDefaultOwner() {
	var owner source.User
	result := DB.Where("role = ?", source.RoleOwner).First(&owner)

	if result.Error == gorm.ErrRecordNotFound {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("owner123"), bcrypt.DefaultCost)
		owner = source.User{
			Email:    "owner@shop.com",
			Password: string(hashedPassword),
			Username: "Владелец",
			Role:     source.RoleOwner,
			IsActive: true,
		}

		if err := DB.Create(&owner).Error; err != nil {
			log.Printf("Failed to create default owner: %v", err)
		} else {
			log.Println("Default owner created: owner@shop.com / owner123")
		}
	}
}

func CreteBD() {
	os.Create("store.db")
}

func InitDB() {
	CreteBD()
	ConnectDB()
	MigrateDB()
	createDefaultOwner()
	InsertSampleData()
}
