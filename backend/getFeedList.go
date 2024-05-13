package backend

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

func GetFeedList() []FeedsInfo {
	result := []FeedsInfo{}

	if dbFilePath == "" {
		log.Fatal("Database file path is not set")
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Print the feeds in the Feeds table
	rows, err := db.Query("SELECT [Link], [Category] FROM [Feeds]")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var link string
		var category string
		err = rows.Scan(&link, &category)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, FeedsInfo{Link: link, Category: category})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}
