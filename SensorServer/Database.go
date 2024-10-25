package main

import (
	"database/sql"
	"fmt"
	"sync"
)
import _ "github.com/mattn/go-sqlite3"

type Database struct {
	db  *sql.DB
	mut sync.Mutex
}

func (self *Database) save_measurements(humidity int, temperature int) error {
	self.mut.Lock()
	defer self.mut.Unlock()
	insertSQL := `INSERT INTO sensor(humidity, temperature) VALUES (?, ?)`
	statement, err := self.db.Prepare(insertSQL)
	if err != nil {
		return err
	}

	defer statement.Close()
	_, err = statement.Exec(humidity, temperature)
	if err != nil {
		return err
	}
	fmt.Println("data inserted:", humidity, temperature)
	return nil
}

func (self *Database) start() (e error) {
	if db, err := sql.Open("sqlite3", "../../dbNew/sensor.db"); err == nil {
		self.db = db
		println("Database opened")
	} else {
		e = err
	}
	return
}

func (self *Database) stop() (e error) {
	if self.db != nil {
		if err := self.db.Close(); err == nil {
			println("Database closed")
		} else {
			e = err
		}
	}
	return
}

func createTable() error { //only for initialize
	// Create a table if it doesn't already exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS sensor (
        Timestamp DATETIME DEFAULT (datetime('now','localtime')),
        humidity INTEGER,
        temperature INTEGER
    );`

	_, err := DB.db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}

var DB = Database{db: nil}
