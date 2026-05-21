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
		{Name: "Ноутбук", Description: "Мощный ноутбук для работы и игр", Price: 59999, Stock: 10, ImageURL: "/images/laptop.jpg"},
		{Name: "Смартфон", Description: "Современный смартфон с отличной камерой", Price: 29999, Stock: 20, ImageURL: "/images/phone.jpg"},
		{Name: "Наушники", Description: "Беспроводные наушники с шумоподавлением", Price: 4999, Stock: 30, ImageURL: "/images/headphones.jpg"},
		{Name: "Клавиатура", Description: "Механическая клавиатура с подсветкой", Price: 3999, Stock: 15, ImageURL: "/images/keyboard.jpg"},
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
