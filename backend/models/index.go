package models

import (
	"bilet/backend/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

var DB *gorm.DB

func InitDB() {
	dsn := strings.Join([]string{
		"host= 127.0.0.1",        //+ os.Getenv("DB_HOST"),
		"user= postgres",         //+ os.Getenv("DB_USER"),
		"password=Anna2109",      // + os.Getenv("DB_PASSWORD"),
		"dbname=ProfBilet",       // + os.Getenv("DB_NAME"),
		"port= 5432",             //+ os.Getenv("DB_PORT"),
		"sslmode=disable",        //+ os.Getenv("DB_SSL_MODE"),
		"TimeZone=Europe/Moscow", // + os.Getenv("DB_TIMEZONE"),
	}, " ")
	log.Print(dsn)
	log.Print(os.Getenv("DB_NAME"))
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	MigrateAll(db)
	go TokensRepeatedCleanup(config.TokensCleanupPeriod)
	DB = db
}

func MigrateAll(db *gorm.DB) {
	db.AutoMigrate(
		&Accounts{},
		&Token{},
	)
	log.Print("Migrated database")
}
