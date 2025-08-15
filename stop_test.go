package main

import (
	"reflect"
	"testing"
)

func TestStop(t *testing.T) {
	tests := []struct {
		words       []string
		wantedWords []string
		num         int
	}{
		{
			words:       []string{"the", "kailash", "everywhere"},
			wantedWords: []string{"kailash"},
			num:         1,
		},
		{
			words:       []string{"present", "library", "university", "earth"},
			wantedWords: []string{"library", "university", "earth"},
			num:         2,
		},
		{
			words:       []string{"the", "backward", "before", "earth", "jupiter", "mars"},
			wantedWords: []string{"earth", "jupiter", "mars"},
			num:         3,
		},
	}

	for _, test := range tests {
		filteredWords := stop(test.words)
		if !reflect.DeepEqual(filteredWords, test.wantedWords) {
			t.Errorf("TEST %d FAILED stop = %v, want %v", test.num, filteredWords, test.wantedWords)
		}
	}
}
