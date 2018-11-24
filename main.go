package main

import (
	"database/sql"
	"event-management/handler"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	const connStr = "postgres://qposcndr:5cuTSlv1gdl8KtwqULYYmMDCy-9RsgC4@stampy.db.elephantsql.com:5432/qposcndr"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(handler.StartServer(":8000", db))
}
