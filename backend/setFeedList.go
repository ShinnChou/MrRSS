package backend

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

func SetFeedList(feeds []FeedsInfo) {
	if dbFilePath == "" {
		log.Fatal("Database file path is not set")
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert feeds into the Feeds table
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO [Feeds]([Link], [Category]) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, feed := range feeds {
		_, err = stmt.Exec(feed.Link, feed.Category)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}
