package main

import (
	"flag"
	"fmt"
)

func main() {
	var index Index

	mode := flag.String("index", "db", "Choose inmem or db mode")
	flag.Parse()

	if *mode == "inmem" {
		fmt.Println("Using in-memory index")
		index = WordIndex{
			wordMap: make(map[string]map[string]int), // word --> map(url it was found in --> number of times it was found in that url)
			urlMap:  make(map[string]int),            // url --> total number of words in that url
		}
	} else {
		fmt.Println("Using database index")
		index = Database{db: database()}
	}

	go server(index)
	select {}
}

// https://usf-cs272-s25.github.io/top10/
// http://localhost:8080/top10/
