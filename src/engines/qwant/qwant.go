package qwant

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
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

	col.OnResponse(func(r *colly.Response) {
		var pageStr string = r.Ctx.Get("page")
		if pageStr == "" {
			// If I'm using GET as the first page
			return
		}

		page, _ := strconv.Atoi(pageStr)

		var parsedResponse QwantResponse
		err := json.Unmarshal(r.Body, &parsedResponse)
		if err != nil {
			log.Error().Err(err).Msgf("%v: Failed body unmarshall to json:\n%v", Info.Name, string(r.Body))
		}

		mainline := parsedResponse.Data.Res.Items.Mainline
		counter := 1
		for _, group := range mainline {
			if group.Type != "web" {
				continue
			}
			for _, result := range group.Items {
				goodURL := parse.ParseURL(result.URL)

				res := bucket.MakeSEResult(goodURL, result.Title, result.Description, Info.Name, page, counter)
				bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
				counter += 1
			}
		}
	})

	localeParam := getLocale(&options)
	nRequested := settings.RequestedResultsPerPage
	deviceParam := getDevice(&options)
	safeSearchParam := getSafeSearch(&options)

	for i := 0; i < options.MaxPages; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		reqString := Info.URL + query + "&count=" + strconv.Itoa(nRequested) + localeParam + "&offset=" + strconv.Itoa(i*nRequested) + deviceParam + safeSearchParam

		sedefaults.DoGetRequest(reqString, colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

// qwant returns this array when an invalid locale is supplied
var validLocales = [...]string{"bg_bg", "br_fr", "ca_ad", "ca_es", "ca_fr", "co_fr", "cs_cz", "cy_gb", "da_dk", "de_at", "de_ch", "de_de", "ec_ca", "el_gr", "en_au", "en_ca", "en_gb", "en_ie", "en_my", "en_nz", "en_us", "es_ad", "es_ar", "es_cl", "es_co", "es_es", "es_mx", "es_pe", "et_ee", "eu_es", "eu_fr", "fc_ca", "fi_fi", "fr_ad", "fr_be", "fr_ca", "fr_ch", "fr_fr", "gd_gb", "he_il", "hu_hu", "it_ch", "it_it", "ko_kr", "nb_no", "nl_be", "nl_nl", "pl_pl", "pt_ad", "pt_pt", "ro_ro", "sv_se", "th_th", "zh_cn", "zh_hk"}

func getLocale(options *engines.Options) string {
	locale := strings.ToLower(options.Locale)
	for _, vl := range validLocales {
		if locale == vl {
			return "&locale=" + locale
		}
	}
	log.Warn().Msgf("qwant.getLocale(): Invalid qwant locale (%v) supplied. Defaulting to en_US. Qwant supports these (disregard specific formatting): %v", options.Locale, validLocales)
	return "&locale=" + strings.ToLower(config.DefaultLocale)
}

func getDevice(options *engines.Options) string {
	if options.Mobile {
		return "&device=mobile"
	}
	return "&device=desktop"
}

func getSafeSearch(options *engines.Options) string {
	if options.SafeSearch {
		return "&safesearch=1"
	}
	return "&safesearch=0"
}

/*
col.OnHTML("div[data-testid=\"sectionWeb\"] > div > div", func(e *colly.HTMLElement) {
	//first page
	idx := e.Index

	dom := e.DOM
	baseDOM := dom.Find("div[data-testid=\"webResult\"] > div > div > div > div > div")
	hrefElement := baseDOM.Find("a[data-testid=\"serTitle\"]")
	linkHref, hrefExists := hrefElement.Attr("href")
	linkText := parse.ParseURL(linkHref)
	titleText := strings.TrimSpace(hrefElement.Text())
	descText := strings.TrimSpace(baseDOM.Find("div > span").Text())

	if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
		var pageStr string = e.Request.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, -1, page, idx+1)
		bucket.AddSEResult(res, Info.Name, relay, options, pagesCol)
	} else {
		log.Info().Msgf("Not Good! %v\n%v\n%v", linkText, titleText, descText)
	}
})
*/
