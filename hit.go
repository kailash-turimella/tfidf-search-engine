package main

type Hit struct {
	Title string
	Url   string
	Score float64
}

// Implementing function in sort.Interface
type ByScore []Hit

func (a ByScore) Len() int      { return len(a) }
func (a ByScore) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool {
	if a[i].Score == a[j].Score {
		return a[i].Url > a[j].Url
	}
	return a[i].Score > a[j].Score
}
