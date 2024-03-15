package database

import (
	"fmt"
	"log"
	"os"

	"github.com/antiloger/nhostel-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBintance struct {
	Db *gorm.DB
}

var DB DBintance

func Connect() DBintance {
	fmt.Print("dbf")
	dsn := "host=localhost user=postgres password=7913456 dbname=teledata port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("[error] Failed to connect to database! \n", err)
		os.Exit(2)
	}

	log.Println("[Info] Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	db.AutoMigrate(&models.User{})

	DB = DBintance{
		Db: db,
	}

	return DB
}