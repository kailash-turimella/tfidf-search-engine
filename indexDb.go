package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kljensen/snowball"
)

func (db Database) search(term string) ([]Hit, error) {
	term = strings.ToLower(term)

	term, err := snowball.Stem(term, "english", true)
	if err != nil {
		fmt.Println("error stemming term: ", err)
		return []Hit{}, err
	}

	var word Word
	result := db.db.Where("word = ?", term).First(&word)

	// Check if word exists
	if result.Error != nil {
		fmt.Println(term, " not found")
		return []Hit{}, nil
	}

	rankedUrls := []Hit{}

	var wordCount []WordCount // a slice that contains all rows in the wordCount database which contain the term
	result = db.db.Where("word_id = ?", word.ID).Find(&wordCount)

	if result.Error != nil {
		fmt.Println("Error retrieving data:", result.Error)
		return []Hit{}, result.Error
	}

	for _, wc := range wordCount {

		var url URL
		result = db.db.First(&url, wc.URLID)
		if result.Error != nil {
			fmt.Println("Error retrieving url: ", result.Error)
			return []Hit{}, result.Error
		}

		var numURLs int64
		db.db.Model(&URL{}).Count(&numURLs)

		// numOccurrencesInURL = wc.Count
		// totalNumWordsInURL = url.NumWords
		// totalNumURLs = numURLs
		// numURLsContainingTerm := len(wordCount)

		rankedUrls = append(rankedUrls, Hit{Title: getTitle(url.URL), Url: url.URL, Score: tfIdf(wc.Count, url.NumWords, int(numURLs), len(wordCount))})
	}

	sort.Sort(ByScore(rankedUrls))
	return rankedUrls, nil
}

// NUM URLS
func (db Database) numUrls() int {
	var numURLs int64
	db.db.Model(&URL{}).Count(&numURLs)
	return int(numURLs)
}

// Adds a URL to the index
func (db Database) addUrl(urlStr string) bool {
	var urlObj URL
	result := db.db.FirstOrCreate(&urlObj, URL{URL: urlStr})

	return result.RowsAffected > 0
}

// Populates the index with a list of URLs
func (db Database) populateIndex(urlStr string, words []string) {
	var urlObj URL
	db.db.FirstOrCreate(&urlObj, URL{URL: urlStr})

	numWords := 0                   // number of words in the url
	wordMap := make(map[string]int) // map of words in the url and their frequencies
	// update the database with the words found
	for _, wordStr := range words {
		// add word to wordMap
		wordMap[wordStr] += 1
		// increment wordCount
		numWords += 1
	}

	// add words from the wordMap to the wordBatch and insert wordBatch into the database
	var wordBatch []Word
	for wordStr := range wordMap {
		var word Word
		db.db.FirstOrCreate(&word, Word{Word: wordStr})
		wordBatch = append(wordBatch, word)
	}

	// add words from the wordBatch to the wordCountBatch and insert wordCountBatch into the database
	var wordCountBatch []WordCount
	for _, word := range wordBatch {
		wordCountBatch = append(wordCountBatch, WordCount{WordID: word.ID, URLID: urlObj.ID, Count: wordMap[word.Word]})
	}
	db.db.Create(wordCountBatch)

	// update the total number of words in the url
	db.db.Model(&urlObj).Update("NumWords", numWords)
}
