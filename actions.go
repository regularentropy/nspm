package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

/*
NOTICE:
If function returns -1 or nothing, it means that something went wrong inside it
If function returns different values, it means that it works as expected
Why not throwing err ? Isn't simple + more code to handle
*/

/* ============= Actions used for category manipulation ============= */
func createCategory(cats *[]Categories) {
	var cat Categories
	cat.CategoryName = input("Enter category title: ")
	*cats = append(*cats, cat)
}

func editCategory(cats *[]Categories) {
	clearScreen()
	cat_index := chooseCategoryIndex(cats)
	if cat_index >= 0 && cat_index <= len(*cats) {
		category_menu(&(*cats)[cat_index])
	}
}

func renameCategory(cats *[]Categories) {
	if len(*cats) == 0 {
		return
	}
	clearScreen()
	listCategories(cats)
	u_index := input_int(": ")
	if u_index == -1 {
		return
	}
	(*cats)[u_index].CategoryName = input("Enter new title: ")

}

func listCategories(cats *[]Categories) {
	fmt.Println("Available categories: ")
	for index, c := range *cats {
		fmt.Printf("%d : %s\n", index, c.CategoryName)
		for _, r := range c.Records {
			fmt.Printf("\t%s\n", r.Title)
		}
	}
}

func removeCategory(cats *[]Categories) {
	if len(*cats) == 0 {
		return
	}
	clearScreen()
	c_cat := chooseCategoryIndex(cats)
	if c_cat == -1 {
		return
	}
	*cats = append((*cats)[:c_cat], (*cats)[c_cat+1:]...)
}

func chooseCategoryIndex(cats *[]Categories) int {
	if len(*cats) == 0 {
		return -1
	}
	listCategories(cats)
	u_index := input_int(": ")
	if u_index == -1 {
		return -1
	}
	return u_index
}

/* ============= Actions used for record manipulation ============= */

func createRecord(recs *[]Record) {
	var rec Record
	rec.Title = input("Enter site: ")
	rec.Username = input("Enter username: ")
	rec.Password = input("Enter password: ")
	rec.Description = input("Enter description: ")
	*recs = append(*recs, rec)
}

func editRecord(recs *[]Record) {
	clearScreen()
	rec_index := chooseRecordIndex(recs)
	if rec_index >= 0 && rec_index <= len(*recs) {
		record_menu(&(*recs)[rec_index])
	}
}

func removeRecord(recs *[]Record) {
	clearScreen()
	c_rec := chooseRecordIndex(recs)
	if c_rec == -1 {
		return
	}
	*recs = append((*recs)[:c_rec], (*recs)[c_rec+1:]...)
}

func listRecords(recs *[]Record) {
	fmt.Println("Available records:")
	for index, r := range *recs {
		fmt.Printf("%d: %s\n", index, r.Title)
		fmt.Printf("\tusername: %s\n", r.Username)
		fmt.Printf("\tpassword: %s\n", r.Password)
		fmt.Printf("\tdescription: %s\n", r.Description)
	}
}

func chooseRecordIndex(recs *[]Record) int {
	if len(*recs) == 0 {
		return -1
	}
	listRecords(recs)
	u_index := input_int(": ")
	if u_index == -1 {
		return -1
	}
	return u_index
}

/* Move the record between categories */
func moveRecord(cats *[]Categories) {
	var selected_rec Record
	var destin_category *Categories
	var source_category *Categories
	var cat_index_to int
	var cat_index_from int
	var rec_index_from int
	if len(*cats) == 0 {
		return
	}
	clearScreen()
	fmt.Print("Select category to move from\n\n")
	cat_index_from = chooseCategoryIndex(cats)
	if cat_index_from >= 0 && cat_index_from <= len(*cats) {
		clearScreen()
		fmt.Print("Select record to move\n\n")
		rec_index_from = chooseRecordIndex(&(*cats)[cat_index_from].Records)
		source_category = &(*cats)[cat_index_from]
		if rec_index_from >= 0 && rec_index_from <= len(source_category.Records) {
			selected_rec = source_category.Records[rec_index_from]
		} else {
			return
		}
	} else {
		return
	}
	clearScreen()
	fmt.Print("Select destination category\n\n")
	cat_index_to = chooseCategoryIndex(cats)
	if cat_index_to >= 0 && cat_index_to <= len(*cats) {
		destin_category = &(*cats)[cat_index_to]
		destin_category.Records = append(destin_category.Records, selected_rec)
		source_category.Records = append(source_category.Records[:rec_index_from], source_category.Records[rec_index_from+1:]...)
	} else {
		return
	}
}

/* ============= Actions for records editing ============= */
func changeRecordTitle(title *string) {
	fmt.Printf("Current title: %s\n", *title)
	*title = input("Enter new title: ")
}

func changeRecordUsername(username *string) {
	fmt.Printf("Current username: %s\n", *username)
	*username = input("Enter new username: ")
}

func changeRecordPassword(password *string) {
	fmt.Printf("Current password: %s\n", *password)
	*password = input("Enter new password: ")
}

func changeRecordDescription(description *string) {
	fmt.Printf("Current description: %s\n", *description)
	*description = input("Enter new description: ")
}

func generateRecordPassword(password *string) {
	ps_length := input_int("Enter length: ")
	if ps_length > 0 {
		ps := make([]rune, ps_length)
		for i := range ps {
			ps[i] = rand.Int31n(126-33) + 33
		}
		*password = string(ps)
	}
}

func listRecord(rec *Record) {
	clearScreen()
	fmt.Printf("%s\n", rec.Title)
	fmt.Printf("Username: %s\n", rec.Username)
	fmt.Printf("Password: %s\n", rec.Password)
	fmt.Printf("Description: %s\n", rec.Description)
}

/* ============= Actions responsible for manipulation with the database ============= */

/* Creates a new encrypted database */
func createNewDatabase() {
	init_rec := &Categories{} /* A dummy record. Must be to init the database */

	fmt.Println("[Creating a new database]")

	db_name := input("Database name: ")
	if db_name == "" {
		fmt.Println("Database title can't be empty")
		os.Exit(0)
	}
	db_key_plain, err := createNewPassword()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	db_key := getDerivedPassword(&db_key_plain)
	db_path := filepath.Join(getDatabaseFolder(), db_name)
	test_rec, _ := json.Marshal(init_rec)
	encrypt(&test_rec, &db_path, db_key)
}

/* Return marshalled database */
func marshalDatabase(cats *[]Categories) *[]byte {
	var m_db []byte
	m_db, _ = json.Marshal(*cats)
	return &m_db
}

/* Return unmarshalled database  */
func unmarshalDatabase(db []byte) *[]Categories {
	var u_db []Categories
	json.Unmarshal(db, &u_db)
	return &u_db
}

/* Getting the location of the folder with the databases */
func getDatabaseFolder() string {
	u_dir, _ := os.UserHomeDir()
	fpath := filepath.Join(u_dir, ".nanopm")
	os.MkdirAll(fpath, 0700)
	return fpath
}

/* Return full path to the selected database */
func getDatabasePath(databaseFolder string) (string, string) {
	var u_index int
	files, _ := os.ReadDir(databaseFolder)
	if len(files) > 0 {
		for {
			fmt.Println("Select database to open")
			for index, file := range files {
				fmt.Printf("%d : %s\n", index, file.Name())
			}
			u_index = input_int(": ")
			if u_index > len(files) || u_index < 0 {
				fmt.Println("Database doesn't exist")
			} else {
				break
			}
		}
		filepath := filepath.Join(databaseFolder, files[u_index].Name())
		return files[u_index].Name(), filepath
	}
	return "", ""
}

/* ============= Misc actions ============= */

/* Replacement of the default input function. Allows entering more than one word */
func input(text string) string {
	fmt.Print(text)
	inputReader := bufio.NewReader(os.Stdin)
	uInput, _ := inputReader.ReadString('\n')
	return strings.TrimSpace(uInput)
}

/* Replacement for the default int input function. Checks if the user didn't enter anything and returns -1 in that case */
func input_int(text string) int {
	var num int
	fmt.Print(text)
	_, err := fmt.Scanf("%d", &num)
	if err != nil {
		return -1
	}
	return num
}

/* Function responsible for entering a password when opening a database */
func readDatabasePassword() *[]byte {
	fmt.Printf("Enter password: ")
	password, _ := term.ReadPassword(0)
	return &password
}

/* Clearing screen after some action */
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

/* Enter to continue field */
func enterToContinue() {
	fmt.Println("\nPress ENTER to continue")
	fmt.Scanln()
}
