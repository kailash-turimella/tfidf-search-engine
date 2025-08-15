package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func extract(body []byte) ([]string, []string) {

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error parsing HTML")
		return nil, nil
	}

	var words []string
	var hrefs []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
			words = append(words, strings.Fields(n.Data)...)
		}

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					hrefs = append(hrefs, attr.Val)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return words, hrefs
}
