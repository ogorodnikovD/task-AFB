package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ogorodnikovD/task-AFB/internal/processor"
	"github.com/ogorodnikovD/task-AFB/utils"
)

func (app *application) uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.ServerError(w, r, err)
		return
	}
	defer file.Close()

	fmt.Printf("Filename: %s\n", header.Filename)
	fmt.Printf("Size: %d bytes\n", header.Size)

	records, err := processor.ReadCSV(w, r, file)
	if err != nil {
		utils.ServerError(w, r, err)
	}
	clientsList, err := processor.ProcessCSVData(records)
	if err != nil {
		utils.ServerError(w, r, err)
	}

	err = processor.InsertClients(app.clients.DB, clientsList)
	if err != nil {
		utils.ServerError(w, r, err)
	}

	fmt.Println("Successful data recording")
	os.Exit(0)
}
