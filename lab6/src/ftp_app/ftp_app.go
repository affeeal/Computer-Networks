package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_HOST = "students.yss.su"
	DB_NAME = "iu9networkslabs"
	DB_LOGIN = "iu9networkslabs"
	DB_PASSWORD = "Je2dTYr6"

	PATH = "../../media/"

	ADDRESS = "students.yss.su:21"
	PASSWORD = "3Ru7yOTA"
	USER = "ftpiu8"
)

func main() {
	log.Printf("Connecting to %s database...\n", DB_NAME)
	db, err := sql.Open("mysql", DB_LOGIN + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ")/" + DB_NAME)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	log.Println("Querying the data...")
	rows, err := db.Query("SELECT title FROM iu9afanasyev")
	if err != nil {
		log.Fatal("Failed to query: ", err)
	}
	defer rows.Close()

	log.Println("Creating a file...")
	name := "Ilya-Afanasyev_" + time.Now().Format("2006-01-02_15-04-05") + ".md"
	file, err := os.Create(PATH + name)
	if err != nil {
		log.Fatal("Failed to create a file: ", err)
	}
	defer file.Close()

	log.Println("Writing the data...")
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			log.Fatal("Failed to scan: ", err)
		}

		if _, err := file.WriteString("- " + title + ";\n"); err != nil {
			log.Fatal("Failed to write a title: ", err)
		}
	}

	log.Printf("Connecting to %s FTP-server...\n", ADDRESS)
	conn, err := ftp.Dial(ADDRESS, ftp.DialWithTimeout(5 * time.Second))
	if err != nil {
	    log.Fatal("Failed to dial: ", err)
	}

	err = conn.Login(USER, PASSWORD)
	if err != nil {
	    log.Fatal("Failed to login: ", err)
	}

	readFile, err := os.Open(PATH + name)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer readFile.Close()

	log.Printf("Storing %s...\n", name)
	if err = conn.Stor(name, readFile); err != nil {
		log.Fatal("Failed to issue STOR FTP command: ", err)
	}

	if err := conn.Quit(); err != nil {
	    log.Fatal("Failed to quit: ", err)
	}

	log.Println("Completed successfully.")
}
