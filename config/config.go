package config

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
    "log"
    "os"
)

var DB *gorm.DB

func ConnectDatabase() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Retrieve the database user and password from environment variables
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")

    // Construct DSN (Data Source Name) for MySQL connection
    dsn := dbUser + ":" + dbPassword + "@tcp(mysql:3306)/questions_db?charset=utf8mb4&parseTime=True&loc=Local"

    // Open a connection to the MySQL database
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }

    DB = db
}