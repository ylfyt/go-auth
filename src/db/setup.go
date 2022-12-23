package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-auth/src/config"
)

const MAX_DB_CONN = 50

var dbConnChan = make(chan *sql.DB, MAX_DB_CONN)

func init() {
	for i := 0; i < MAX_DB_CONN; i++ {
		dbConnChan <- nil
	}
}

func BorrowDbConnection() (*sql.DB, error) {
	conn := <-dbConnChan
	// Connect to DB
	newDB, err := sql.Open("postgres", config.DB_CONNECTION)
	if err != nil {
		dbConnChan <- conn
		return nil, err
	}

	err = newDB.Ping()
	if err != nil {
		dbConnChan <- conn
		return nil, err
	}
	return newDB, nil
}

func ReturnDbConnection(dbConn *sql.DB) {
	// Close connection
	err := dbConn.Close()
	if err != nil {
		fmt.Println("Failed to close connection")
	}
	dbConnChan <- dbConn
}
