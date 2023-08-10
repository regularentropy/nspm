package main

import "fmt"

func main_menu(db *[]Categories, db_name *string, db_path *string, db_key *[]byte) {
	in_menu := true
	for in_menu {
		clearScreen()
		fmt.Printf("nspm v%d.%d.%d\n", major, minor, patch)
		fmt.Printf("Editing [%s]\n", *db_name)
		fmt.Printf(
			"1.Create new category\n" +
				"2.Edit category\n" +
				"3.Rename category\n" +
				"4.List categories\n" +
				"5.Move records\n" +
				"6.Remove category\n" +
				"7.Change database password\n" +
				"8.Save and exit\n")
		u_choice := input_int(": ")
		switch u_choice {
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
			changeDatabasePassword(db_key)
		case 8:
			encrypt(marshalDatabase(db), db_path, db_key)
			in_menu = false
		}
	}
}

/* Menu for editing the records */
func category_menu(cat *Categories) {
	in_menu := true
	for in_menu {
		clearScreen()
		fmt.Printf("[Editing '%s']\n", cat.CategoryName)
		fmt.Printf(
			"1.Add new record\n" +
				"2.Edit record\n" +
				"3.Remove record\n" +
				"4.List records\n" +
				"5.To main menu\n")
		u_choice := input_int(": ")
		switch u_choice {
		case 1:
			createRecord(&cat.Records)
		case 2:
			editRecord(&cat.Records)
		case 3:
			removeRecord(&cat.Records)
		case 4:
			clearScreen()
			listRecords(&cat.Records)
			enterToContinue()
		case 5:
			in_menu = false
		}
	}
}

func record_menu(rec *Record) {
	in_menu := true
	for in_menu {
		clearScreen()
		fmt.Printf("Editing %s\n", rec.Title)
		fmt.Printf(
			"1.Change title\n" +
				"2.Change username\n" +
				"3.Change password\n" +
				"4.Change description\n" +
				"5.Generate password\n" +
				"6.List current record\n" +
				"7.To record menu\n")
		u_choice := input_int(": ")
		switch u_choice {
		case 1:
			changeRecordField("title", &rec.Title)
		case 2:
			changeRecordField("username", &rec.Username)
		case 3:
			changeRecordField("password", &rec.Password)
		case 4:
			changeRecordField("description", &rec.Title)
		case 5:
			generateRecordPassword(&rec.Password)
		case 6:
			clearScreen()
			listRecord(rec)
			enterToContinue()
		case 7:
			in_menu = false
		}
	}
}
