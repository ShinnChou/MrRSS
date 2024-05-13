package backend

import (
	"database/sql"
	"log"
	"sort"

	_ "github.com/glebarez/go-sqlite"
)

func GetHistory() []FeedContentsInfo {
	result := []FeedContentsInfo{}

	if dbFilePath == "" {
		log.Fatal("Database file path is not set")
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the items in the History table
	rows, err := db.Query("SELECT [FeedTitle], [FeedImage], [Title], [Link], [TimeSince], [Time], [Image], [Content], [Readed] FROM [History]")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var feedTitle string
		var feedImage string
		var title string
		var link string
		var timeSince string
		var time string
		var image string
		var content string
		var readed bool
		err = rows.Scan(&feedTitle, &feedImage, &title, &link, &timeSince, &time, &image, &content, &readed)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, FeedContentsInfo{FeedTitle: feedTitle, FeedImage: feedImage, Title: title, Link: link, TimeSince: timeSince, Time: time, Image: image, Content: content, Readed: readed})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// Sort the result by time
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time > result[j].Time
	})

	return result
}

func CheckInHistory(feed FeedContentsInfo) bool {
	if dbFilePath == "" {
		log.Fatal("Database file path is not set")
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the item is in the History table
	rows, err := db.Query("SELECT [Link] FROM [History] WHERE [Link] = ?", feed.Link)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return rows.Next()
}

func GetHistoryReaded(feed FeedContentsInfo) bool {
	if dbFilePath == "" {
		log.Fatal("Database file path is not set")
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the Readed field in the History table
	rows, err := db.Query("SELECT [Readed] FROM [History] WHERE [Link] = ?", feed.Link)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var readed bool
	if rows.Next() {
		err = rows.Scan(&readed)
		if err != nil {
			log.Fatal(err)
		}
	}

	return readed
}
