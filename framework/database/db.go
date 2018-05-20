package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func InitDB(dsn string) {
	d, err := sql.Open("postgres", dsn)
	if nil != err {
		log.Fatal("cannot open db connection", dsn, err)
	}
	db = d
}

func Get() *sql.DB {
	if nil == db {
		log.Fatal("db is not initialized")
	}
	return db
}
