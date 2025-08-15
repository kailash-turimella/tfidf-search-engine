package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/kljensen/snowball"
)

type PageData struct {
	urlStr string
	words  []string
	hrefs  []string
}

func crawl(startURL string, index Index) int {
	var wg sync.WaitGroup
	rules := robots(startURL)

	crawlQueue := make(chan string, 10000)
	downloadQueue := make(chan string, 10000)
	indexQueue := make(chan PageData, 10000)

	go func() {
		for urlStr := range crawlQueue {
			_, err := url.ParseRequestURI(urlStr)
			if err != nil {
				continue
			}
			if !isAllowed(urlStr, rules) {
				wg.Done()
				continue
			}
			if !index.addUrl(urlStr) {
				wg.Done()
				continue
			}
			downloadQueue <- urlStr
		}
	}()

	// Start download workers
	for range 5 {
		go func() {
			for urlStr := range downloadQueue {
				time.Sleep(time.Duration(rules.CrawlDelay) * time.Second)

				body, err := download(urlStr)
				if err != nil {
					wg.Done()
					continue
				}
				words, hrefs := extract(body)

				indexQueue <- PageData{
					urlStr: urlStr,
					words:  words,
					hrefs:  clean(urlStr, hrefs),
				}
			}
		}()
	}

	go func() {
		for page := range indexQueue {
			words := stop(page.words)
			for i, word := range words {
				stemmed, err := snowball.Stem(word, "english", true)
				if err == nil {
					words[i] = stemmed
				}
			}
			index.populateIndex(page.urlStr, words)

			for _, href := range page.hrefs {
				hostUrl, err := url.Parse(page.urlStr)
				if err != nil {
					fmt.Println("error parsing url: ", err)
					continue
				}
				hrefUrl, err := url.Parse(href)
				if err != nil {
					fmt.Println("error parsing url: ", err)
					continue
				}
				if hrefUrl.Host != hostUrl.Host {
					continue
				}
				crawlQueue <- href
				wg.Add(1)
			}
			wg.Done()
		}
	}()

	// Start crawling
	crawlQueue <- startURL
	wg.Add(1)

	wg.Wait()

	// Close all channels
	close(crawlQueue)
	close(downloadQueue)
	close(indexQueue)

	return index.numUrls()
}
