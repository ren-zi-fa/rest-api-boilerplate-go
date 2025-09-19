package main

import (
	"api-todolist/cmd/api"
	"api-todolist/config"
	"api-todolist/db"
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {

	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		DBName:               config.Envs.DBName,
		Addr:                 config.Envs.DBAddress,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	initStorage(db)
	if err != nil {
		log.Fatal("error connect to db", err)
	}

	defer db.Close()
	server := api.NewServerAPI(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal("error", err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal("error ping to db", err)
	}
	log.Println("success connect to db")
}
