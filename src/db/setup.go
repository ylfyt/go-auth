package db

import "database/sql"

const MAX_DB_CONN = 50

type DbConnection struct {
	Id int
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

func BorrowDbConnection() *DbConnection {
	conn := <-dbConnChan
	// Connect to DB
	return conn
}

func ReturnDbConnection(dbConn *DbConnection) {
	// Close connection
	dbConnChan <- dbConn
}
