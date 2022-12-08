package db

import (
	"database/sql"
	"fmt"
	"go-auth/src/config"
	_ "github.com/lib/pq"
)

const MAX_DB_CONN = 50

type DbConnection struct {
	Id    int
	SqlDb *sql.DB
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
	newConn, err := sql.Open("postgres", config.DB_CONNECTION)
	if err != nil {
		dbConnChan <- conn
		return nil, err
	}
	conn.SqlDb = newConn
	return conn, nil
}

func ReturnDbConnection(dbConn *DbConnection) {
	// Close connection
	err := dbConn.SqlDb.Close()
	if err != nil {
		fmt.Println("Failed to close connection")
	}
	dbConnChan <- dbConn
}
