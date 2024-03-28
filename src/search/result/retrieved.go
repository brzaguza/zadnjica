package result

import "github.com/hearchco/hearchco/src/search/engines"

// variables are 1-indexed
// Information about what Rank a result was on some Search Engine
type RetrievedRank struct {
	SearchEngine engines.Name
	Rank         uint
	Page         uint
	OnPageRank   uint
}

// The info a Search Engine returned about some Result
type RetrievedResult struct {
	URL         string
	URLHash     string
	Title       string
	Description string
	ImageResult ImageResult
	Rank        RetrievedRank
}
