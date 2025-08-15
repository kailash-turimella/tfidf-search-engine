package main

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCrawl(t *testing.T) {
	tests := []struct {
		url         string
		wantedLinks int
		num         int
	}{
		// Basic case
		{
			url:         "https://usf-cs272-s25.github.io/test-data/lab02a/index.html",
			wantedLinks: 4,
			num:         1,
		},
		{
			url:         "https://usf-cs272-s25.github.io/test-data/lab02a/href.html",
			wantedLinks: 2,
			num:         2,
		},
		{
			url:         "https://usf-cs272-s25.github.io/test-data/lab02a/simple.html",
			wantedLinks: 1,
			num:         3,
		},
		{
			url:         "https://usf-cs272-s25.github.io/test-data/lab02a/style.html",
			wantedLinks: 3,
			num:         4,
		},
	}

	for _, test := range tests {

		// setup and populate database
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			fmt.Println("Error creating the database: ", err)
		}
		// Delete current tables
		db.Migrator().DropTable("words")
		db.Migrator().DropTable("urls")
		db.Migrator().DropTable("word_counts")

		// Create new tables
		db.Exec("PRAGMA foreign_keys = ON;")
		err = db.AutoMigrate(&Word{}, &URL{}, &WordCount{})
		if err != nil {
			fmt.Println("Error creating tables: ", err)
		}

		var dbIndex Index
		dbIndex = Database{db: db}
		numLinksVisitedDb := crawl(test.url, dbIndex)

		var inmemIndex Index
		inmemIndex = WordIndex{
			wordMap: make(map[string]map[string]int),
			urlMap:  make(map[string]int),
		}
		numLinksVisitedInmem := crawl(test.url, inmemIndex)

		if numLinksVisitedInmem != numLinksVisitedDb || numLinksVisitedDb != test.wantedLinks {
			t.Errorf("\nurl: %s\nTEST %d FAILED got db:%v inmem:%v links, wanted %v links", test.url, test.num, numLinksVisitedDb, numLinksVisitedInmem, test.wantedLinks)
		}
	}
}
