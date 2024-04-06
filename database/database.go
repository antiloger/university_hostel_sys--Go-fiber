package database

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antiloger/nhostel-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBintance struct {
	Db *gorm.DB
}

var DB DBintance

func Connect() {
	// Create a new reader, assuming input will be provided via the standard input device (os.Stdin)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter database host:")
	host, _ := reader.ReadString('\n')

	fmt.Println("Enter database user:")
	user, _ := reader.ReadString('\n')

	fmt.Println("Enter database password:")
	password, _ := reader.ReadString('\n')

	fmt.Println("Enter database name:")
	dbname, _ := reader.ReadString('\n')

	fmt.Println("Enter database port:")
	port, _ := reader.ReadString('\n')

	fmt.Println("Enter database ssl mode:")
	sslmode, _ := reader.ReadString('\n')

	fmt.Println("Enter database TimeZone:")
	timeZone, _ := reader.ReadString('\n')

	// Remove the newline character from the input strings
	host = strings.TrimSpace(host)
	user = strings.TrimSpace(user)
	password = strings.TrimSpace(password)
	dbname = strings.TrimSpace(dbname)
	port = strings.TrimSpace(port)
	sslmode = strings.TrimSpace(sslmode)
	timeZone = strings.TrimSpace(timeZone)

	// Construct the DSN using formatted strings
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("[error] Failed to connect to database! \n", err)
		os.Exit(2)
	}

	log.Println("[Info] Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	db.AutoMigrate(&models.UserInfo{}, &models.Hostel{}, &models.Student{}, &models.HostelOwner{}, &models.Booking{}, &models.Admin{}, &models.Article{}, &models.Warden{})

	DB = DBintance{
		Db: db,
	}
}
