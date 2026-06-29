package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	// Ganti dengan user desa dan password desa123
	dsn := "desa:desa123@tcp(127.0.0.1:3306)/dbdesa?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal konek database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	fmt.Println("✅ Koneksi database berhasil!")
}