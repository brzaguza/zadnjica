package startpage

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/bucket"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/hearchco/hearchco/src/sedefaults"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, &options, &timings)

	sedefaults.PagesColRequest(Info.Name, pagesCol, ctx)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, ctx)
	sedefaults.ColError(Info.Name, col)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, hrefExists := dom.Find(dompaths.Link).Attr("href")
		linkText := parse.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
			page := sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			pageRankCounter[page]++
		} else {
			log.Trace().
				Str("engine", Info.Name.String()).
				Str("url", linkText).
				Str("title", titleText).
				Str("description", descText).
				Msg("Matched result, but couldn't retrieve data")
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "to prevent possible abuse of our service") {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("Request blocked by engine due to scraping")
		} else if strings.Contains(string(r.Body), "This page cannot function without javascript") {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("Engine couldn't load requests, needs javascript")
		}
	})

	safeSearch := getSafeSeach(&options)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query + safeSearch
	anonUrll := Info.URL + anonymize.String(query) + safeSearch
	sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, &retError)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + "&page=" + strconv.Itoa(i+1) + safeSearch
		anonUrll := Info.URL + anonymize.String(query) + "&page=" + strconv.Itoa(i+1) + safeSearch
		sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getSafeSeach(options *engines.Options) string {
	if options.SafeSearch {
		return "" // for startpage, Safe Search is the default
	}
	return "&qadf=none"
}
