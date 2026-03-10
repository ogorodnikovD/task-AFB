package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ogorodnikovD/task-AFB/internal/processor"
	"github.com/ogorodnikovD/task-AFB/storage"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	clients *processor.ClientModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "root:rP7i8kje4@/bank_clients?parseTime=true", "MySQL data source name")
	flag.Parse()

	db, err := storage.OpenDB(*dsn)
	if err != nil {
		log.Println("DB opening error", err)
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		clients: &processor.ClientModel{DB: db},
	}

	log.Println("Starting server on ", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	if err != nil {
		log.Fatal("Starting server error: ", err)
		os.Exit(1)
	}
}
