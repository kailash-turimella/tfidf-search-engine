package main

import (
	"slices"
	"testing"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		body                     []byte
		wantedWords, wantedHrefs []string
		num                      int
	}{
		// Basic case
		{
			body:        []byte(`<html><body><h1>Hello World!</h1><p>This is a test.</p><a href="https://example.com">here</a></body></html>`),
			wantedWords: []string{"Hello", "World!", "This", "is", "a", "test.", "here"},
			wantedHrefs: []string{"https://example.com"},
			num:         1,
		},
		{
			body:        []byte(``),
			wantedWords: []string{},
			wantedHrefs: []string{},
			num:         2,
		},
		{
			body:        []byte(`<html><body><p>No links here!</p></body></html>`),
			wantedWords: []string{"No", "links", "here!"},
			wantedHrefs: []string{},
			num:         3,
		},
		{
			body:        []byte(`<html><body><a href="https://example.com"></a></body></html>`),
			wantedWords: []string{},
			wantedHrefs: []string{"https://example.com"},
			num:         4,
		},
		{
			body:        []byte(`<html><body><a href="https://example.com">First</a> <a href="https://example.com">Second</a></body></html>`),
			wantedWords: []string{"First", "Second"},
			wantedHrefs: []string{"https://example.com", "https://example.com"},
			num:         5,
		},
		{
			body:        []byte(`<html><body><div><p>Hello <strong>bold</strong> text.</p></div><a href="https://example.com">Click</a></body></html>`),
			wantedWords: []string{"Hello", "bold", "text.", "Click"},
			wantedHrefs: []string{"https://example.com"},
			num:         6,
		},
		{
			body:        []byte(`<html><body><p>Broken <img src="image.jpg"/> paragraph.</p><a href="https://example.com"/></body></html>`),
			wantedWords: []string{"Broken", "paragraph."},
			wantedHrefs: []string{"https://example.com"},
			num:         7,
		},
		{
			body:        []byte(`<html><body><p>Special &amp; characters &lt;test&gt;</p><a href="https://example.com">here</a></body></html>`),
			wantedWords: []string{"Special", "&", "characters", "<test>", "here"},
			wantedHrefs: []string{"https://example.com"},
			num:         8,
		},
		{
			body:        []byte(`<html><body><a href="https://example.com"></a><a href="https://example.com">Valid</a></body></html>`),
			wantedWords: []string{"Valid"},
			wantedHrefs: []string{"https://example.com", "https://example.com"},
			num:         9,
		},
		{
			body:        []byte(`<html><body><h1>Header</h1><p>First paragraph.</p><p>Second paragraph with <a href="https://example.com">link</a></p></body></html>`),
			wantedWords: []string{"Header", "First", "paragraph.", "Second", "paragraph", "with", "link"},
			wantedHrefs: []string{"https://example.com"},
			num:         10,
		},
	}

	for _, test := range tests {
		words, hrefs := extract(test.body)
		if !slices.Equal(words, test.wantedWords) {
			t.Errorf("TEST %d FAILED words = %v, want %v", test.num, words, test.wantedWords)
		}
		if !slices.Equal(hrefs, test.wantedHrefs) {
			t.Errorf("TEST %d FAILED hrefs = %v, want %v", test.body, hrefs, test.wantedHrefs)
		}
	}
}
