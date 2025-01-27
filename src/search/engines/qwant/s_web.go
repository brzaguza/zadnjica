package qwant

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper/parse"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/moreurls"
)

func (se Engine) WebSearch(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)

	se.OnResponse(func(r *colly.Response) {
		var pageStr string = r.Ctx.Get("page")
		if pageStr == "" {
			// If I'm using GET as the first page
			return
		}

		pageIndex := se.PageFromContext(r.Request.Ctx)
		page := pageIndex + opts.Pages.Start + 1

		var parsedResponse jsonResponse
		if err := json.Unmarshal(r.Body, &parsedResponse); err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Bytes("body", r.Body).
				Msg("Failed to parse response, couldn't unmarshal JSON")
		}

		mainline := parsedResponse.Data.Res.Items.Mainline
		counter := 1
		for _, group := range mainline {
			if group.Type != "web" {
				continue
			}
			for _, jsonResult := range group.Items {
				goodURL, goodTitle, goodDesc := parse.SanitizeFields(jsonResult.URL, jsonResult.Title, jsonResult.Description)

				r, err := result.ConstructResult(se.Name, goodURL, goodTitle, goodDesc, page, counter)
				if err != nil {
					log.Error().
						Caller().
						Err(err).
						Str("url", goodURL).
						Str("title", goodTitle).
						Str("desc", goodDesc).
						Int("page", page).
						Int("rank", counter).
						Msg("Failed to construct result")
				} else {
					log.Trace().
						Caller().
						Int("page", page).
						Int("rank", counter).
						Str("result", fmt.Sprintf("%v", r)).
						Msg("Sending result to channel")
					resChan <- r
					counter++
				}
			}
		}
	})

	// Constant params.
	paramLocaleV := localeParamValue(opts.Locale)
	paramSafeSearchV := safeSearchParamValue(opts.SafeSearch)

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Build the parameters.
		params := moreurls.NewParams(
			paramQueryK, query,
			paramCountK, paramCountV,
			paramLocaleK, paramLocaleV,
			paramPageK, strconv.Itoa(pageNum0*10),
			paramSafeSearchK, paramSafeSearchV,
		)

		// Build the url.
		urll := moreurls.Build(searchURL, params)

		// Build anonymous url, by anonymizing the query.
		params.Set(paramQueryK, anonymize.String(query))
		anonUrll := moreurls.Build(searchURL, params)

		// Send the request.
		if err := se.Get(ctx, urll, anonUrll); err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
