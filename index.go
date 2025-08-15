package main

import (
	"gorm.io/gorm"
)

type Index interface {
	search(string) ([]Hit, error)
	numUrls() int
	addUrl(string) bool
	populateIndex(string, []string)
}

type WordIndex struct {
	wordMap map[string]map[string]int // word --> map(url it was found in --> number of times it was found in that url)
	urlMap  map[string]int            // url --> total number of words in that url
}

type Database struct {
	db *gorm.DB
}
