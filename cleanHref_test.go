package main

import (
	"slices"
	"testing"
)

func TestCleanHref(t *testing.T) {
	tests := []struct {
		num         int
		host        string
		hrefs, want []string
	}{
		{
			1, "https://cs272-f24.github.io/",
			[]string{"/", "/help/", "/syllabus/", "https://gobyexample.com/"},
			[]string{"https://cs272-f24.github.io/", "https://cs272-f24.github.io/help/", "https://cs272-f24.github.io/syllabus/", "https://gobyexample.com/"},
		},
		{
			2, "https://example.com",
			[]string{"http://example.com/page1", "https://example.com/page2", "/page3"},
			[]string{"http://example.com/page1", "https://example.com/page2", "https://example.com/page3"},
		},
		{
			3, "https://blog.example.com",
			[]string{"/about", "https://example.com/home"},
			[]string{"https://blog.example.com/about", "https://example.com/home"},
		},
		{
			4, "https://mywebsite.com",
			[]string{"about", "/contact", "https://other.com/page"},
			[]string{"https://mywebsite.com/about", "https://mywebsite.com/contact", "https://other.com/page"},
		},
	}

	for _, test := range tests {
		got := clean(test.host, test.hrefs)
		if !slices.Equal(got, test.want) {
			t.Errorf("TEST %d FAILED: got %v, expected %v", test.num, got, test.want)
		}
	}
}
