package connections

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewConnection() (*sql.DB,error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_school")
	if err != nil {
		return nil, err
	}
	
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db,nil
}
