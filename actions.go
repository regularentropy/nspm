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
func createCategory(categories *[]Category) {
	var category Category
	category.CategoryName = input("Enter the category title: ")
	*categories = append(*categories, category)
}

func editCategory(categories *[]Category) {
	clearScreen()
	categoryIndex := chooseCategoryIndex(categories)
	if categoryIndex >= 0 && categoryIndex < len(*categories) {
		categoryMenu(&(*categories)[categoryIndex])
	}
}

func renameCategory(categories *[]Category) {
	if len(*categories) == 0 {
		return
	}
	clearScreen()
	listCategories(categories)
	userIndex := inputInt(": ")
	if userIndex == -1 {
		return
	}
	(*categories)[userIndex].CategoryName = input("Enter the new title: ")
}

func listCategories(categories *[]Category) {
	fmt.Println("Available categories:")
	for index, category := range *categories {
		fmt.Printf("%d : %s\n", index, category.CategoryName)
		for _, record := range category.Records {
			fmt.Printf("\t%s\n", record.Title)
		}
	}
}

func removeCategory(categories *[]Category) {
	if len(*categories) == 0 {
		return
	}
	clearScreen()
	categoryIndex := chooseCategoryIndex(categories)
	if categoryIndex == -1 {
		return
	}
	*categories = append((*categories)[:categoryIndex], (*categories)[categoryIndex+1:]...)
}

func chooseCategoryIndex(categories *[]Category) int {
	if len(*categories) == 0 {
		return -1
	}
	listCategories(categories)
	userIndex := inputInt(": ")
	if userIndex == -1 {
		return -1
	}
	return userIndex
}

/* ============= Actions used for record manipulation ============= */

func createRecord(records *[]Record) {
	var record Record
	record.Title = input("Enter the site name: ")
	record.Username = input("Enter the username: ")
	record.Password = input("Enter the password: ")
	record.Description = input("Enter the description: ")
	*records = append(*records, record)
}

func editRecord(records *[]Record) {
	clearScreen()
	recordIndex := chooseRecordIndex(records)
	if recordIndex >= 0 && recordIndex < len(*records) {
		recordMenu(&(*records)[recordIndex])
	}
}

func removeRecord(records *[]Record) {
	clearScreen()
	recordIndex := chooseRecordIndex(records)
	if recordIndex == -1 {
		return
	}
	*records = append((*records)[:recordIndex], (*records)[recordIndex+1:]...)
}

func listRecords(records *[]Record) {
	fmt.Println("Available records:")
	for index, record := range *records {
		fmt.Printf("%d: %s\n", index, record.Title)
		fmt.Printf("\tUsername: %s\n", record.Username)
		fmt.Printf("\tPassword: %s\n", record.Password)
		fmt.Printf("\tDescription: %s\n", record.Description)
	}
}

func chooseRecordIndex(records *[]Record) int {
	if len(*records) == 0 {
		return -1
	}
	listRecords(records)
	userIndex := inputInt(": ")
	if userIndex == -1 {
		return -1
	}
	return userIndex
}

/* Move the record between categories */
func moveRecord(categories *[]Category) {
	clearScreen()
	fmt.Println("Select the category to move from:")
	categoryIndexFrom := chooseCategoryIndex(categories)
	if categoryIndexFrom < 0 || categoryIndexFrom >= len(*categories) {
		return
	}

	sourceCategory := &(*categories)[categoryIndexFrom]
	clearScreen()
	fmt.Println("Select the record to move:")
	recordIndexFrom := chooseRecordIndex(&sourceCategory.Records)
	if recordIndexFrom < 0 || recordIndexFrom >= len(sourceCategory.Records) {
		return
	}

	selectedRecord := sourceCategory.Records[recordIndexFrom]
	clearScreen()
	fmt.Println("Select the destination category:")
	categoryIndexTo := chooseCategoryIndex(categories)
	if categoryIndexTo < 0 || categoryIndexTo >= len(*categories) {
		return
	}

	destinationCategory := &(*categories)[categoryIndexTo]
	destinationCategory.Records = append(destinationCategory.Records, selectedRecord)
	sourceCategory.Records = append(sourceCategory.Records[:recordIndexFrom], sourceCategory.Records[recordIndexFrom+1:]...)
}

/* ============= Actions for records editing ============= */

func changeRecordField(fieldName string, field *string) {
	fmt.Printf("Current %s: %s\n", fieldName, *field)
	*field = input(fmt.Sprintf("Enter new %s: ", fieldName))
}

func generateRecordPassword(password *string) {
	passwordLength := inputInt("Enter the length: ")
	if passwordLength > 0 {
		passwordRunes := make([]rune, passwordLength)
		for i := range passwordRunes {
			passwordRunes[i] = rand.Int31n(126-33) + 33
		}
		*password = string(passwordRunes)
	}
}

func displayRecord(record *Record) {
	clearScreen()
	fmt.Println(record.Title)
	fmt.Printf("Username: %s\n", record.Username)
	fmt.Printf("Password: %s\n", record.Password)
	fmt.Printf("Description: %s\n", record.Description)
}

/* ============= Actions responsible for manipulation with the database ============= */

/* Creates a new encrypted database */
func createNewDatabase() {
	fmt.Println("[Creating a new database]")

	databaseName := input("Enter the database name: ")
	if databaseName == "" {
		fmt.Println("The database name cannot be empty")
		os.Exit(0)
	}

	databaseKeyPlain, err := createNewPassword()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	databaseKey := getDerivedPassword(&databaseKeyPlain)

	databasePath := filepath.Join(getDatabaseFolder(), databaseName)
	initRec := &Category{} // A dummy record to initialize the database
	testRec, _ := json.Marshal(initRec)
	encrypt(&testRec, &databasePath, databaseKey)
}

/* Return marshaled database */
func marshalDatabase(categories *[]Category) *[]byte {
	marshaledDatabase, _ := json.Marshal(*categories)
	return &marshaledDatabase
}

/* Return unmarshaled database  */
func unmarshalDatabase(database []byte) *[]Category {
	var unmarshaledDatabase []Category
	json.Unmarshal(database, &unmarshaledDatabase)
	return &unmarshaledDatabase
}

/* Getting the location of the folder with the databases */
func getDatabaseFolder() string {
	userDir, _ := os.UserHomeDir()
	folderPath := filepath.Join(userDir, ".nspm")
	if err := os.MkdirAll(folderPath, 0700); err != nil {
		log.Fatal(err)
	}
	return folderPath
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
		userIndex := inputInt(": ")
		if userIndex >= 0 && userIndex < len(files) {
			filePath := filepath.Join(databaseFolder, files[userIndex].Name())
			return files[userIndex].Name(), filePath
		}
		fmt.Println("Database doesn't exist")
	}
}

/* ============= Misc actions ============= */

/* Replacement of the default input function. Allows entering more than one word */
func input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func inputInt(prompt string) int {
	fmt.Print(prompt)
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
