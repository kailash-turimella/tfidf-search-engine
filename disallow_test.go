package main

import (
	"testing"
)

func TestDisallow(t *testing.T) {
	tests := []struct {
		url         string
		wantedLinks int
		num         int
	}{
		{
			url:         "http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/index.html",
			wantedLinks: 27,
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
	}

	for _, test := range tests {
		index := WordIndex{
			wordMap: make(map[string]map[string]int),
			urlMap:  make(map[string]int),
		}

		numUrls := crawl(test.url, index)

		if numUrls != test.wantedLinks {
			t.Errorf("Expected %d links, got %d\n", test.wantedLinks, numUrls)
		}
	}
}
