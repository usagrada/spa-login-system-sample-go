package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// DB is the database connection
var DB *sql.DB

// Connect connects to the database
func Setup() {
	fmt.Println("Connecting to database...")
	mysqlConfig := mysql.Config{
		DBName:    "login_system",
		User:      "root",
		Passwd:    "",
		Addr:      "localhost:3306",
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
	}
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database!")
	DB = db
}

func Initialize() {
	// Create tables
	createTables()
}

func createTables() {
	// Create users table
	_, err := DB.Exec("CREATE TABLE IF NOT EXISTS users (id int NOT NULL AUTO_INCREMENT, username varchar(255) NOT NULL, password varchar(255) NOT NULL, PRIMARY KEY (id))")
	if err != nil {
		panic(err)
	}
}
