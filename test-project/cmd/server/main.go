package main

import (
	"database/sql"
	"net/http"
	"os"
)

// Hardcoded credentials - should be detected
const apiKey = "sk-1234567890abcdef"

func getUserByID(db *sql.DB, userID string) error {
	// SQL injection vulnerability
	query := "SELECT * FROM users WHERE id = " + userID
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func readConfig() {
	// Missing error handling
	data, _ := os.ReadFile("config.json")
	println(string(data))
}

func main() {
	http.ListenAndServe(":8080", nil)
}
