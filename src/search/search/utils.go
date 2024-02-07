package search

import (
	"strings"

	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func procBang(query *string, options *engines.Options, conf *config.Config) (config.Timings, []engines.Name) {
	useSpec, specEng := procSpecificEngine(*query, options, conf)
	goodCat := procCategory(*query, options)
	if !goodCat && !useSpec && (*query)[0] == '!' {
		// options.category is set to GENERAL
		log.Debug().
			Str("query", *query).
			Msg("search.procBang(): invalid bang (not category or engine shortcut)")
	}

	trimBang(query)

	if useSpec {
		return conf.Categories[category.GENERAL].Timings, []engines.Name{specEng}
	} else {
		return conf.Categories[options.Category].Timings, conf.Categories[options.Category].Engines
	}
}

func trimBang(query *string) {
	if (*query)[0] == '!' {
		*query = strings.SplitN(*query, " ", 2)[1]
	}
}

func procSpecificEngine(query string, options *engines.Options, conf *config.Config) (bool, engines.Name) {
	if query[0] != '!' {
		return false, engines.UNDEFINED
	}
	sp := strings.SplitN(query, " ", 2)
	bangWord := sp[0][1:]
	for key, val := range conf.Settings {
		if strings.EqualFold(bangWord, val.Shortcut) || strings.EqualFold(bangWord, key.String()) {
			return true, key
		}
	}

	return false, engines.UNDEFINED
}

func procCategory(query string, options *engines.Options) bool {
	cat := category.FromQuery(query)
	if cat != "" {
		options.Category = cat
	}
	if options.Category == "" {
		options.Category = category.GENERAL
	}
	return cat != ""
}
