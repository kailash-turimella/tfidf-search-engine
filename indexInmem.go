package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kljensen/snowball"
)

func (indexMap WordIndex) search(term string) ([]Hit, error) {
	term = strings.ToLower(term)

	term, err := snowball.Stem(term, "english", true)
	if err != nil {
		fmt.Println("error stemming term: ", err)
		return []Hit{}, err
	}

	termMap, exists := indexMap.wordMap[term]
	if !exists {
		return []Hit{}, err
	}

	rankedUrls := []Hit{}

	numUrls := len(termMap) // number of urls containing the term

	// Makes sure we didn't get an empty map
	if numUrls == 0 {
		return rankedUrls, nil
	}

	for url, wordCount := range termMap {
		// make sure the url has at least one word
		if indexMap.urlMap[url] == 0 {
			continue
		}

		// add a Hit with the url and the corresponding tf-idf score to the slice of ranked urls
		rankedUrls = append(rankedUrls, Hit{Title: getTitle(url), Url: url, Score: tfIdf(wordCount, indexMap.urlMap[url], len(indexMap.urlMap), len(termMap))})
	}
	sort.Sort(ByScore(rankedUrls))

	return rankedUrls, nil
}

// Returns the number of urls in the index
func (indexMap WordIndex) numUrls() int {
	return len(indexMap.urlMap)
}

// Adds a URL to the index
func (indexMap WordIndex) addUrl(urlStr string) bool {
	// ignore the url if it's already been visited
	if _, exists := indexMap.urlMap[urlStr]; exists {
		return false
	}
	indexMap.urlMap[urlStr] = 0 // number of words in the url
	return true
}

// Populates the index with a list of URLs
func (indexMap WordIndex) populateIndex(urlStr string, words []string) {
	// update the urlMap and populate indexMap
	for _, word := range words {
		indexMap.urlMap[urlStr] += 1
		word = strings.ToLower(word)
		if _, exists := indexMap.wordMap[word]; !exists {
			indexMap.wordMap[word] = make(map[string]int)
			indexMap.wordMap[word][urlStr] = 0
		}
		indexMap.wordMap[word][urlStr]++
	}
}
