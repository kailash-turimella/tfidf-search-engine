package main

import (
	"testing"
	"time"
)

func TestCrawlDelay(t *testing.T) {
	t1 := time.Now()

	var index Index
	index = WordIndex{
		wordMap: make(map[string]map[string]int),
		urlMap:  make(map[string]int),
	}

	crawl("http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/index.html", index)
	t2 := time.Now()
	if t2.Sub(t1) < (10 * time.Second) {
		t.Errorf("TestDisallow was too fast\n")
	}
}
