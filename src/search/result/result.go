package result

import (
	"github.com/gocolly/colly/v2"
)

type ImageFormat struct {
	Height uint
	Width  uint
}

type ImageResult struct {
	Original         ImageFormat
	Thumbnail        ImageFormat
	ThumbnailURL     string
	ThumbnailURLHash string
	Source           string
	SourceURL        string
}

// Everything about some Result, calculated and compiled from multiple search engines
// The URL is the primary key
type Result struct {
	URL         string
	URLHash     string
	Rank        uint
	Score       float64
	Title       string
	Description string
	EngineRanks []RetrievedRank
	ImageResult ImageResult
	Response    *colly.Response
}
