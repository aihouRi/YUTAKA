package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB // 声明一个全局的数据库句柄变量

func InitDB(dsn string) error {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	fmt.Println("sql.Open() executed successfully. db handle created.")

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database (ping failed): %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	return nil
}