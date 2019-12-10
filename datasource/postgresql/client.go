package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	// import PostgreSQL driver
	_ "github.com/lib/pq"
)

// Init connects to the database server
func Init(address, username, password string, dbName string) (db *sqlx.DB, err error) {

	db, err = sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%v/%s", username, password, address, dbName))
	if err == nil {
		err = db.Ping()
	}

	return
}
