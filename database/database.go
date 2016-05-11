package database

import (
	"database/sql"
	"fmt"
	// Add mysql support for sql package
	_ "github.com/go-sql-driver/mysql"
)

var (
	// Connection stores the AAC MySQL connection
	Connection *sql.DB
)

// NewConnection opens a new connection to the MySQL server
func NewConnection(user, password, database string) error {
	c, err := sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", user, password, database))
	if err != nil {
		return err
	}
	Connection = c
	return nil
}
