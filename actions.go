package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/term"
)

/*
NOTICE:
If the function returns -1 or nothing, it means that something went wrong inside it.
If the function returns a different value, it means that it works as expected.
Why not throw an error? It's simpler and requires less code to handle.
*/

/* ============= Actions used for category manipulation ============= */
func createCategory(cats *[]Categories) {
	var cat Categories
	cat.CategoryName = input("Enter the category title: ")
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
	(*cats)[u_index].CategoryName = input("Enter the new title: ")
}

func listCategories(cats *[]Categories) {
	fmt.Println("Available categories:")
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
	rec.Title = input("Enter the site name: ")
	rec.Username = input("Enter the username: ")
	rec.Password = input("Enter the password: ")
	rec.Description = input("Enter the description: ")
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
		fmt.Printf("\tUsername: %s\n", r.Username)
		fmt.Printf("\tPassword: %s\n", r.Password)
		fmt.Printf("\tDescription: %s\n", r.Description)
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
	clearScreen()
	fmt.Println("Select the category to move from:")
	catIndexFrom := chooseCategoryIndex(cats)
	if catIndexFrom < 0 || catIndexFrom >= len(*cats) {
		return
	}

	sourceCategory := &(*cats)[catIndexFrom]
	clearScreen()
	fmt.Println("Select the record to move:")
	recIndexFrom := chooseRecordIndex(&sourceCategory.Records)
	if recIndexFrom < 0 || recIndexFrom >= len(sourceCategory.Records) {
		return
	}

	selectedRec := sourceCategory.Records[recIndexFrom]
	clearScreen()
	fmt.Println("Select the destination category:")
	catIndexTo := chooseCategoryIndex(cats)
	if catIndexTo < 0 || catIndexTo >= len(*cats) {
		return
	}

	destinCategory := &(*cats)[catIndexTo]
	destinCategory.Records = append(destinCategory.Records, selectedRec)
	sourceCategory.Records = append(sourceCategory.Records[:recIndexFrom], sourceCategory.Records[recIndexFrom+1:]...)
}

/* ============= Actions for records editing ============= */

func changeRecordField(fieldName string, field *string) {
	fmt.Printf("Current %s: %s\n", fieldName, *field)
	*field = input(fmt.Sprintf("Enter new %s: ", fieldName))
}

func generateRecordPassword(password *string) {
	psLength := input_int("Enter the length: ")
	if psLength > 0 {
		ps := make([]rune, psLength)
		for i := range ps {
			ps[i] = rand.Int31n(126-33) + 33
		}
		*password = string(ps)
	}
}

func listRecord(rec *Record) {
	clearScreen()
	fmt.Println(rec.Title)
	fmt.Printf("Username: %s\n", rec.Username)
	fmt.Printf("Password: %s\n", rec.Password)
	fmt.Printf("Description: %s\n", rec.Description)
}

/* ============= Actions responsible for manipulation with the database ============= */

/* Creates a new encrypted database */
func createNewDatabase() {
	fmt.Println("[Creating a new database]")

	db_name := input("Enter the database name: ")
	if db_name == "" {
		fmt.Println("The database name cannot be empty")
		os.Exit(0)
	}

	db_key_plain, err := createNewPassword()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	db_key := getDerivedPassword(&db_key_plain)

	db_path := filepath.Join(getDatabaseFolder(), db_name)
	init_rec := &Categories{} // A dummy record. Must be to init the database
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
	uDir, _ := os.UserHomeDir()
	fPath := filepath.Join(uDir, ".nspm")
	if err := os.MkdirAll(fPath, 0700); err != nil {
		log.Fatal(err)
	}
	return fPath
}

/* Return full path to the selected database */
func getDatabasePath(databaseFolder string) (string, string) {
	files, _ := os.ReadDir(databaseFolder)
	if len(files) == 0 {
		return "", ""
	}

	for {
		fmt.Println("Select a database to open:")
		for index, file := range files {
			fmt.Printf("%d : %s\n", index, file.Name())
		}
		uIndex := input_int(": ")
		if uIndex >= 0 && uIndex < len(files) {
			filepath := filepath.Join(databaseFolder, files[uIndex].Name())
			return files[uIndex].Name(), filepath
		}
		fmt.Println("Database doesn't exist")
	}
}

/* ============= Misc actions ============= */

/* Replacement of the default input function. Allows entering more than one word */
func input(text string) string {
	fmt.Print(text)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func input_int(text string) int {
	fmt.Print(text)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	num, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	return num
}

/* Function responsible for entering a password when opening a database */
func readDatabasePassword() *[]byte {
	fmt.Print("Enter the password: ")
	password, _ := term.ReadPassword(int(os.Stdin.Fd()))
	return &password
}

/* Clearing screen after some action */
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else { // Linux or Mac
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

/* Enter to continue field */
func enterToContinue() {
	fmt.Println("\nPress ENTER to continue")
	fmt.Scanln()
}
