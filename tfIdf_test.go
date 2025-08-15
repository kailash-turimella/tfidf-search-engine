package main

import (
	"math"
	"testing"
)

func TestTfIdf(t *testing.T) {
	tests := []struct {
		numOccurrencesInURL int
		wordsInURL          int
		numURLs             int
		numURLsWithTerm     int
		expected            float64
		num                 int
	}{
		{
			numOccurrencesInURL: 3,
			wordsInURL:          10,
			numURLs:             3,
			numURLsWithTerm:     1,
			expected:            (3.0 / 10.0) * math.Log(3.0/2.0),
			num:                 1,
		},
		{
			numOccurrencesInURL: 0,
			wordsInURL:          10,
			numURLs:             2,
			numURLsWithTerm:     0,
			expected:            0,
			num:                 2,
		},
		{
			numOccurrencesInURL: 10,
			wordsInURL:          10,
			numURLs:             3,
			numURLsWithTerm:     1,
			expected:            math.Log(3.0 / 2.0),
			num:                 3,
		},
	}

	for _, test := range tests {
		result := tfIdf(test.numOccurrencesInURL, test.wordsInURL, test.numURLs, test.numURLsWithTerm)

		if result != test.expected {
			t.Errorf("TEST %d FAILED: got %v, want %v", test.num, result, test.expected)
		}
	}
}
