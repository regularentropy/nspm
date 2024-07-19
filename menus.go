package main

import "fmt"

func mainMenu(db *[]Category, dbName *string, dbPath *string, dbKey *[]byte) {
	inMenu := true
	for inMenu {
		clearScreen()
		fmt.Printf("nspm v%d.%d.%d\n", major, minor, patch)
		fmt.Printf("Editing [%s]\n", *dbName)
		fmt.Printf(
			"1. Create new category\n" +
				"2. Edit category\n" +
				"3. Rename category\n" +
				"4. List categories\n" +
				"5. Move records\n" +
				"6. Remove category\n" +
				"7. Change database password\n" +
				"8. Save and exit\n")
		userChoice := inputInt(": ")
		switch userChoice {
		case 1:
			createCategory(db)
		case 2:
			editCategory(db)
		case 3:
			renameCategory(db)
		case 4:
			clearScreen()
			listCategories(db)
			enterToContinue()
		case 5:
			moveRecord(db)
		case 6:
			removeCategory(db)
		case 7:
			changeDatabasePassword(dbKey)
		case 8:
			encrypt(marshalDatabase(db), dbPath, dbKey)
			inMenu = false
		}
	}
}

/* Menu for editing the records */
func categoryMenu(category *Category) {
	inMenu := true
	for inMenu {
		clearScreen()
		fmt.Printf("[Editing '%s']\n", category.CategoryName)
		fmt.Printf(
			"1. Add new record\n" +
				"2. Edit record\n" +
				"3. Remove record\n" +
				"4. List records\n" +
				"5. To main menu\n")
		userChoice := inputInt(": ")
		switch userChoice {
		case 1:
			createRecord(&category.Records)
		case 2:
			editRecord(&category.Records)
		case 3:
			removeRecord(&category.Records)
		case 4:
			clearScreen()
			listRecords(&category.Records)
			enterToContinue()
		case 5:
			inMenu = false
		}
	}
}

func recordMenu(record *Record) {
	inMenu := true
	for inMenu {
		clearScreen()
		fmt.Printf("Editing %s\n", record.Title)
		fmt.Printf(
			"1. Change title\n" +
				"2. Change username\n" +
				"3. Change password\n" +
				"4. Change description\n" +
				"5. Generate password\n" +
				"6. List current record\n" +
				"7. To record menu\n")
		userChoice := inputInt(": ")
		switch userChoice {
		case 1:
			changeRecordField("title", &record.Title)
		case 2:
			changeRecordField("username", &record.Username)
		case 3:
			changeRecordField("password", &record.Password)
		case 4:
			changeRecordField("description", &record.Description)
		case 5:
			generateRecordPassword(&record.Password)
		case 6:
			clearScreen()
			displayRecord(record)
			enterToContinue()
		case 7:
			inMenu = false
		}
	}
}
