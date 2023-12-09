package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/hearchco/src/bucket/result"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search"
)

func Search(c *gin.Context, config *config.Config, db cache.DB) error {
	var query, pages, deepSearch string

	switch c.Request.Method {
	case "", "GET":
		{
			query = c.Query("q")
			pages = c.DefaultQuery("pages", "1")
			deepSearch = c.DefaultQuery("deep", "false")
		}
	case "POST":
		{
			query = c.PostForm("q")
			pages = c.DefaultPostForm("pages", "1")
			deepSearch = c.DefaultPostForm("deep", "false")
		}
	}

	if query == "" {
		c.String(http.StatusOK, "")
	} else {
		maxPages, err := strconv.Atoi(pages)
		if err != nil {
			log.Debug().Err(err).Msgf("router.Search(): cannot convert \"%v\" to int, reverting to default value of 1", pages)
			maxPages = 1
		}

		visitPages := false
		if deepSearch != "false" {
			log.Trace().Msgf("doing a deep search because deep is: %v", deepSearch)
			visitPages = true
		}

		options := engines.Options{
			MaxPages:   maxPages,
			VisitPages: visitPages,
		}

		var results []result.Result
		gerr := db.Get(query, &results)
		if gerr != nil {
			return fmt.Errorf("router.Search(): failed accessing cache for query %v. error: %w", query, gerr)
		}
		foundInDB := results != nil

		if foundInDB {
			log.Debug().Msgf("Found results for query (%v) in cache", query)
		} else {
			log.Debug().Msg("Nothing found in cache, doing a clean search")

			results = search.PerformSearch(query, options, config)
		}

		resultsShort := result.Shorten(results)
		if resultsJson, err := json.Marshal(resultsShort); err != nil {
			c.String(http.StatusInternalServerError, "")
			return fmt.Errorf("router.Search(): failed marshalling results: %v\n with error: %w", resultsShort, err)
		} else {
			c.String(http.StatusOK, string(resultsJson))
		}

		if !foundInDB {
			serr := db.Set(query, results)
			if serr != nil {
				log.Error().Err(serr).Msgf("router.Search(): error updating database with search results")
			}
		}
	}
	return nil
}
