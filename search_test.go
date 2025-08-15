package main

import (
	"fmt"
	"reflect"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		term     string
		urlMap   map[string]int
		indexMap map[string]map[string]int
		expected []Hit
		num      int
	}{
		// IN-MEMORY TESTING
		{
			term: "word",
			urlMap: map[string]int{
				"https://website.com":  50,
				"https://website2.com": 50,
			},
			indexMap: map[string]map[string]int{
				"word": {
					"https://website.com":  10,
					"https://website2.com": 5,
				},
				"test": {
					"https://website.com":  40,
					"https://website2.com": 45,
				},
			},
			expected: []Hit{
				{"Domain Details Page", "https://website2.com", -0.04054651081081645},
				{"Just a moment...", "https://website.com", -0.0810930216216329},
			},
			num: 1,
		},
		{
			term: "kailash",
			urlMap: map[string]int{
				"https://website.com": 100,
			},
			indexMap: map[string]map[string]int{
				"hello": {
					"https://website.com": 100,
				},
			},
			expected: []Hit{},
			num:      2,
		},

		// DATABASE TESTING
		{
			term:     "romeo",
			urlMap:   map[string]int{},
			indexMap: map[string]map[string]int{},
			expected: []Hit{
				{"Just a moment...", "https://website.com", 0.13862943611198905},
			},
			num: 3,
		},
		{
			term:     "juliet",
			urlMap:   map[string]int{},
			indexMap: map[string]map[string]int{},
			expected: []Hit{
				{"Domain Details Page", "https://website2.com", 0.13862943611198905},
			},
			num: 4,
		},
		{
			term:     "kailash",
			urlMap:   map[string]int{},
			indexMap: map[string]map[string]int{},
			expected: []Hit{},
			num:      5,
		},
	}

	// IN-MEMORY TESTING
	for _, test := range tests {
		if test.num > 2 { // skip DATABASE TESTING
			break
		}
		var index Index
		index = WordIndex{
			wordMap: test.indexMap,
			urlMap:  test.urlMap,
		}

		results, err := index.search(test.term)
		if err != nil {
			t.Errorf("TEST %d FAILED: unexpected error %v", test.num, err)
		}

		if !reflect.DeepEqual(results, test.expected) {
			t.Errorf("TEST %d FAILED: got %v, want %v", test.num, results, test.expected)
		}
	}

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

	urls := []URL{
		{URL: "https://website.com", NumWords: 100},
		{URL: "https://website2.com", NumWords: 100},
		{URL: "https://website3.com", NumWords: 50},
		{URL: "https://website4.com", NumWords: 50},
	}
	words := []Word{
		{Word: "word"},
		{Word: "test"},
		{Word: "romeo"},
		{Word: "juliet"},
	}
	wordCounts := []WordCount{
		{WordID: 1, URLID: 1, Count: 10},
		{WordID: 1, URLID: 2, Count: 5},
		{WordID: 2, URLID: 1, Count: 40},
		{WordID: 2, URLID: 2, Count: 45},
		{WordID: 3, URLID: 1, Count: 20},
		{WordID: 4, URLID: 2, Count: 20},
	}
	db.Create(&urls)
	db.Create(&words)
	db.Create(&wordCounts)

	var dbIndex Index
	dbIndex = Database{db: db}

	for _, test := range tests {
		if test.num <= 2 { // skip IN-MEMORY TESTING
			continue
		}
		results, err := dbIndex.search(test.term)
		if err != nil {
			t.Errorf("TEST %d FAILED: unexpected error %v", test.num, err)
		}
		if !reflect.DeepEqual(results, test.expected) {
			t.Errorf("TEST %d FAILED: got %v, want %v", test.num, results, test.expected)
		}
	}
}
