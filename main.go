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
	var db_name string
	var db_path string
	new_db := flag.Bool("n", false, "Create new database")
	arg_path := flag.String("f", "", "Pipe database into nanopm")
	flag.Parse()

	if *new_db {
		createNewDatabase()
		os.Exit(0)
	}

	if *arg_path != "" {
		db_path = *arg_path
		if _, err := os.Stat(db_path); err != nil {
			fmt.Println("Database doesn't exist")
			os.Exit(0)
		}
		db_name = filepath.Base(db_path)
	}
	if !*new_db && len(*arg_path) == 0 {
		db_name, db_path = getDatabasePath(getDatabaseFolder())
		if db_name == "" && db_path == "" {
			fmt.Println("Databases weren't found. Consider creating new database via -n")
			os.Exit(0)
		}
	}
	db_pass := getDerivedPassword(readDatabasePassword())
	db := unmarshalDatabase(decrypt(*db_pass, db_path))
	main_menu(db, &db_name, &db_path, db_pass)
	clearScreen()
}
