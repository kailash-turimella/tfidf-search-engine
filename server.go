package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func server(index Index) {

	// crawl handler
	http.HandleFunc("/crawl", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query().Get("q")
		if query == "" {
			// http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			http.ServeFile(w, r, "./crawl/index.html")
			return
		}

		// If there is a query, crawl it
		fmt.Println("Crawling URL...")
		linksCrawled := crawl(query, index)
		fmt.Println(linksCrawled, " URLs crawled")
		fmt.Fprintf(w, "Crawled %d links\n\n", linksCrawled)
	})

	// search handler
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query().Get("q")
		if query == "" {
			// http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			http.ServeFile(w, r, "./search/index.html")
			return
		}

		// If there is a query, search it and get a slice of urls sorted based on TfIdf
		results, err := index.search(query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error processing search: %v", err), http.StatusInternalServerError)
			return
		}

		// Check and print if there were any references to the query
		if len(results) == 0 {
			fmt.Fprintf(w, "No results found for query: %s", query)
			fmt.Printf("No results found for query: %s\n", query)
			return
		}

		// Print the results
		fmt.Printf("Found %d result(s) for query: %s\n", len(results), query)
		for _, hit := range results {
			fmt.Fprintf(w, "<a href='%s' target='_blank'>%s</a>  -  Score: %f<br>", hit.Url, hit.Title, hit.Score)
		}
	})

	// Show homepage at "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html_files/index.html")
	})

	// Show robots.txt
	http.Handle("/robots.txt", http.FileServer(http.Dir("./html_files")))

	// Show search page at "/search"
	http.Handle("/search/", http.StripPrefix("/search/", http.FileServer(http.Dir("./html_files/search"))))

	// Show search page at "/crawl"
	http.Handle("/crawl/", http.StripPrefix("/crawl/", http.FileServer(http.Dir("./html_files/crawl"))))

	// Show files from the "top10" directory
	http.Handle("/top10/", http.StripPrefix("/top10/", http.FileServer(http.Dir("./html_files/top10"))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func getTitle(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return url
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	for {
		tokenType := z.Next()
		switch tokenType {
		case html.ErrorToken:
			return url
		case html.StartTagToken:
			token := z.Token()
			if token.Data == "title" {
				z.Next()
				return strings.TrimSpace(z.Token().Data)
			}
		}
	}
}
