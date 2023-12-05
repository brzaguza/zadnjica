package etools

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/hearchco/hearchco/src/bucket"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/hearchco/hearchco/src/sedefaults"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, &options, &timings)

	sedefaults.PagesColRequest(Info.Name, pagesCol, ctx, &retError)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, &ctx, &retError)
	sedefaults.ColError(Info.Name, col, &retError)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkEl := dom.Find(dompaths.Link)
		linkHref, _ := linkEl.Attr("href")
		var linkText string

		if linkHref[0] == 'h' {
			//normal link
			linkText = parse.ParseURL(linkHref)
		} else {
			//telemetry link, e.g. //web.search.ch/r/redirect?event=website&origin=result!u377d618861533351/https://de.wikipedia.org/wiki/Charles_Paul_Wilp
			linkText = parse.ParseURL("http" + strings.Split(linkHref, "http")[1]) //works for https, dont worry
		}

		titleText := strings.TrimSpace(linkEl.Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			pageRankCounter[page]++
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Sorry for the CAPTCHA") {
			log.Error().Msgf("%v: Returned captcha.", Info.Name)
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	err := col.Request("POST", Info.URL, strings.NewReader("query="+query+"&country=web&language=all"), colCtx, nil)
	if engines.IsTimeoutError(err) {
		log.Trace().Err(err).Msgf("%v: failed requesting with POST method", Info.Name)
	} else if err != nil {
		log.Error().Err(err).Msgf("%v: failed requesting with POST method", Info.Name)
	}
	col.Wait() //wait so I can get the JSESSION cookie back

	for i := 1; i < options.MaxPages; i++ {
		pageStr := strconv.Itoa(i + 1)
		colCtx = colly.NewContext()
		colCtx.Put("page", pageStr)

		err := col.Request("GET", pageURL+pageStr, nil, colCtx, nil)
		if engines.IsTimeoutError(err) {
			log.Trace().Err(err).Msgf("%v: failed requesting with GET method on page", Info.Name)
		} else if err != nil {
			log.Error().Err(err).Msgf("%v: failed requesting with GET method on page", Info.Name)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
