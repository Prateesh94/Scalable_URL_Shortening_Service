package database

import (
	"context"
	"fmt"
	"log"
	"main/global"
	"time"
)

var db Database

func init() {
	var err error
	db, err = initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	//	defer db.Close()
}
func CheckUrl(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, _ := db.GetWriteConnection(ctx, "index")
	var exists bool
	err := con.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM index WHERE url = $1)", url).Scan(&exists)
	if err != nil {
		//fmt.Printf("Error checking URL existence: %v", err)
		return false
	}
	if exists {
		//	log.Println("URL already exists in the database.")
		return true

	} else {
		//log.Println("URL does not exist in the database.")
		return false
	}
}
func updateIndex(short_url string) bool {
	attempt := 0
retry:
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, _ := db.GetWriteConnection(ctx, "index")
	_, err := con.Exec(ctx, "INSERT INTO index (url,count) VALUES ($1,$2)", short_url, 0)
	if err != nil && attempt < 3 {

		attempt++
		goto retry

		//fmt.Printf("Error updating index: %v", err)
	} else if attempt >= 3 {
		return false
	}
	return true
}
func UpdateCount(short_url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, _ := db.GetWriteConnection(ctx, "index")
	_, err := con.Exec(ctx, "UPDATE index SET count = count + 1 WHERE url = $1", short_url)
	if err != nil {
		//fmt.Printf("Error updating count: %v", err)
		return false
	}
	return true
}
func deleteIndex(short_url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, _ := db.GetWriteConnection(ctx, "index")
	_, err := con.Exec(ctx, "DELETE FROM index WHERE url = $1", short_url)
	if err != nil {
		//fmt.Printf("Error deleting from index: %v", err)
		return false
	}
	return true
}
func AddData(short_url, original_url string) (global.Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, _ := db.GetWriteConnection(ctx, short_url)
	var data global.Data
	er := con.QueryRow(ctx, "INSERT INTO urls (short_url,original_url) VALUES ($1, $2) RETURNING id, short_url, original_url", short_url, original_url).Scan(&data.Id, &data.ShortUrl, &data.OriginalUrl)
	if er != nil {
		//fmt.Printf("Error inserting data: %v", er)
		return global.Data{}, fmt.Errorf("failed to insert data: %v", er)
	}
	if !updateIndex(short_url) {
		return global.Data{}, fmt.Errorf("failed to update index for short URL: %s", short_url)
	}

	//con.QueryRow(ctx, "SELECT * FROM urls WHERE short_url = $1 AND original_url = $2", short_url, original_url).Scan(&data.Id, &data.OriginalUrl, &data.ShortUrl)
	return data, nil
}

func GetOriginalUrl(short_url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, _ := db.GetReadConnection(ctx, short_url)
	var original_url string
	err := con.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_url = $1", short_url).Scan(&original_url)
	if err != nil {
		//fmt.Printf("Error retrieving original URL: %v", err)
		return "", fmt.Errorf("failed to retrieve original URL: %v", err)
	}
	UpdateCount(short_url) // Update the count in the database
	return original_url, nil
}

func DeleteShortUrl(short_url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, _ := db.GetWriteConnection(ctx, short_url)
	_, err := con.Exec(ctx, "DELETE FROM urls WHERE short_url = $1", short_url)
	if err != nil {
		//fmt.Printf("Error deleting short URL: %v", err)
		return fmt.Errorf("failed to delete short URL: %v", err)
	}
	// Optionally, remove from index
	if !deleteIndex(short_url) {
		return fmt.Errorf("failed to update index after deletion of short URL: %s", short_url)
	}
	return nil
}

func GetCount(short_url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, _ := db.GetWriteConnection(ctx, "index")
	var count int
	err := con.QueryRow(ctx, "SELECT count FROM index WHERE url = $1", short_url).Scan(&count)
	if err != nil {
		//fmt.Printf("Error retrieving count: %v", err)
		return 0, fmt.Errorf("failed to retrieve count for short URL: %v", err)
	}
	return count, nil
}
