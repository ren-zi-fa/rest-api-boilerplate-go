package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	fmt.Println(cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
