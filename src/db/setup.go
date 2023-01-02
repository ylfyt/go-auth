package db

import (
	"database/sql"
	"fmt"
	"go-auth/src/config"

	_ "github.com/lib/pq"
)

const maxConnection = 50

var isPoolMode bool = false

var dbConnChan = make(chan *sql.DB, maxConnection)

func init() {
	if isPoolMode {
		fmt.Println("Setting up db with connection pool")
	} else {
		fmt.Println("Setting up db without connection pool")
	}
	for i := 0; i < maxConnection; i++ {
		if !isPoolMode {
			dbConnChan <- nil
			continue
		}
		temp, err := sql.Open("postgres", config.DB_CONNECTION)
		if err != nil {
			fmt.Println("Failed to open connection :", err)
			dbConnChan <- nil
			continue
		}
		err = temp.Ping()
		if err != nil {
			fmt.Println("Failed to open connection :", err)
			dbConnChan <- nil
			continue
		}
		dbConnChan <- temp
	}
}

func BorrowDbConnection() (*sql.DB, error) {
	conn := <-dbConnChan
	// Connect to DB
	if conn != nil {
		err := conn.Ping()
		if err == nil {
			return conn, nil
		}
		dbConnChan <- conn
		return nil, err
	}

	newDB, err := sql.Open("postgres", config.DB_CONNECTION)
	if err != nil {
		dbConnChan <- nil
		return nil, err
	}

	err = newDB.Ping()
	if err != nil {
		dbConnChan <- nil
		return nil, err
	}
	return newDB, nil
}

func ReturnDbConnection(dbConn *sql.DB) {
	// Close connection
	if isPoolMode {
		dbConnChan <- dbConn
		return
	}

	err := dbConn.Close()
	if err != nil {
		fmt.Println("Failed to close connection")
	}
	dbConnChan <- nil
}

type DbConnection struct {
	Woow int
}

func (me DbConnection) Get() interface{} {
	// Borrow
	fmt.Println("Borrow")
	return DbConnection{}
}

func (me DbConnection) Return(interface{}) {
	fmt.Println("Return")
	// Return
}
