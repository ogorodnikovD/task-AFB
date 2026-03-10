package processor

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ClientData struct {
	ID     string
	Name   string
	Income float64
}

type ClientModel struct {
	DB *sql.DB
}

func ReadCSV(w http.ResponseWriter, r *http.Request, file io.Reader) ([][]string, error) {
	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func ProcessCSVData(records [][]string) ([]ClientData, error) {
	var err error
	Clients := make([]ClientData, len(records)-1)

	for i := 1; i < len(records); i++ {
		income, err := strconv.ParseFloat(records[i][2], 64)

		if err != nil {
			return nil, err
		}

		Clients[i-1] = ClientData{
			ID:     records[i][0],
			Name:   records[i][1],
			Income: income,
		}
	}
	return Clients, err
}

func InsertClients(db *sql.DB, clientList []ClientData) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	valueStrings := make([]string, 0, len(clientList))
	valueArgs := make([]interface{}, 0, len(clientList)*3)

	for _, client := range clientList {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, client.ID, client.Name, client.Income)
	}

	query := fmt.Sprintf("INSERT INTO clients (client_id, client_FIO, client_income) VALUES %s",
		strings.Join(valueStrings, ","))

	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
