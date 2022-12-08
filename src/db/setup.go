package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-auth/src/config"
)

const MAX_DB_CONN = 50

type DbConnection struct {
	Id    int
	sqlDb *sql.DB
}

var dbConns = make([]DbConnection, MAX_DB_CONN)
var dbConnChan = make(chan *DbConnection, MAX_DB_CONN)

func init() {
	for i := 0; i < MAX_DB_CONN; i++ {
		dbConns[i].Id = i
		dbConnChan <- &dbConns[i]
	}
}

func BorrowDbConnection() (*DbConnection, error) {
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
	conn.sqlDb = newDB
	return conn, nil
}

func ReturnDbConnection(dbConn *DbConnection) {
	// Close connection
	err := dbConn.sqlDb.Close()
	if err != nil {
		fmt.Println("Failed to close connection")
	}
	dbConnChan <- dbConn
}
