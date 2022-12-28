package expense

import (
	"database/sql"
	"log"
)

type handler struct {
	database *sql.DB
}

func InitDB(db *sql.DB) *handler {

	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	_, err := db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	return &handler{db}
}
