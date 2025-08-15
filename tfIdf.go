package main

import (
	"math"
)

func tfIdf(numOccurrencesInURL int, wordsInURL int, numURLs int, numURLsWithTerm int) float64 {
	tf := float64(numOccurrencesInURL) / float64(wordsInURL)
	idf := math.Log(float64(numURLs) / (float64(numURLsWithTerm) + 1))
	return tf * idf
}
