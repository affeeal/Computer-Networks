package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mmcdole/gofeed"
)

func clean(db *sql.DB) {
	_, err := db.Exec("TRUNCATE iu9afanasyev")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su)/iu9networkslabs")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	clean(db)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://vesti-k.ru/rss/")
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range feed.Items {
		if _, err = db.Exec(
			"INSERT INTO iu9afanasyev (title, link) VALUES (?, ?)",
			item.Title,
			item.Link,
		); err != nil {
			log.Fatal(err)
		}
	}
}
