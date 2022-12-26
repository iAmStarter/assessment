package expense

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func InitDB() {
	var err error
	url := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

}
