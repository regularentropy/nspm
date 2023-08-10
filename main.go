package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

/* Version */
const (
	major = 1
	minor = 0
	patch = 1
)

/* Example of the record in the database */
type Record struct {
	Title       string `json:"Title"`
	Username    string `json:"Username"`
	Password    string `json:"Password"`
	Description string `json:"Description"`
}

type Categories struct {
	CategoryName string   `json:"Category_Name"`
	Records      []Record `json:"Record"`
}

func main() {
	newDB := flag.Bool("n", false, "Create a new database")
	argPath := flag.String("f", "", "Pipe the database into nanopm")
	flag.Parse()

	var dbName, dbPath string
	var dbPass []byte

	if *newDB {
		createNewDatabase()
		os.Exit(0)
	}

	if *argPath != "" {
		dbPath = *argPath
		if _, err := os.Stat(dbPath); err != nil {
			fmt.Println("The database doesn't exist.")
			os.Exit(0)
		}
		dbName = filepath.Base(dbPath)
	} else {
		dbName, dbPath = getDatabasePath(getDatabaseFolder())
		if dbName == "" && dbPath == "" {
			fmt.Println("Databases were not found. Consider creating a new database via -n")
			os.Exit(0)
		}
	}

	dbPass = *getDerivedPassword(readDatabasePassword())
	db := unmarshalDatabase(decrypt(dbPass, dbPath))
	main_menu(db, &dbName, &dbPath, &dbPass)
	clearScreen()
}
