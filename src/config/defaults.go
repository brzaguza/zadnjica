package config

import (
	"time"

	"github.com/hearchco/hearchco/src/moretime"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

const DefaultLocale string = "en_US"

func EmptyRanking() Ranking {
	rnk := Ranking{
		REXP:    0.5,
		A:       1,
		B:       0,
		C:       1,
		D:       0,
		TRA:     1,
		TRB:     0,
		TRC:     1,
		TRD:     0,
		Engines: map[string]EngineRanking{},
	}

	for _, eng := range engines.Names() {
		rnk.Engines[eng.ToLower()] = EngineRanking{
			Mul:   1,
			Const: 0,
		}
	}

	return rnk
}

func NewRanking() Ranking {
	return EmptyRanking()
}

func NewSettings() map[engines.Name]Settings {
	mp := map[engines.Name]Settings{
		engines.BING: {
			Shortcut: "bi",
		},
		engines.BINGIMAGES: {
			Shortcut: "biimg",
		},
		engines.BRAVE: {
			Shortcut: "br",
		},
		engines.DUCKDUCKGO: {
			Shortcut: "ddg",
		},
		engines.ETOOLS: {
			Shortcut: "ets",
		},
		engines.GOOGLE: {
			Shortcut: "g",
		},
		engines.GOOGLEIMAGES: {
			Shortcut: "gimg",
		},
		engines.GOOGLESCHOLAR: {
			Shortcut: "gs",
		},
		engines.MOJEEK: {
			Shortcut: "mjk",
		},
		engines.PRESEARCH: {
			Shortcut: "ps",
		},
		engines.QWANT: {
			Shortcut: "qw",
		},
		engines.STARTPAGE: {
			Shortcut: "sp",
		},
		engines.SWISSCOWS: {
			Shortcut: "sc",
		},
		engines.YAHOO: {
			Shortcut: "yh",
		},
		engines.YEP: {
			Shortcut: "yep",
		},
	}

	// Check if all search engines have a shortcut set
	for _, eng := range engines.Names() {
		if _, ok := mp[eng]; !ok {
			log.Fatal().
				Str("engine", eng.String()).
				Msg("config.NewSettings(): no shortcut set")
			// ^FATAL
		}
	}

	return mp
}

func NewAllEnabled() []engines.Name {
	return engines.Names()
}

func NewGeneral() []engines.Name {
	return []engines.Name{
		engines.BING,
		engines.BRAVE,
		engines.DUCKDUCKGO,
		engines.ETOOLS,
		engines.GOOGLE,
		engines.MOJEEK,
		engines.PRESEARCH,
		engines.QWANT,
		engines.STARTPAGE,
		engines.SWISSCOWS,
		engines.YAHOO,
		engines.YEP,
	}
}

func NewImage() []engines.Name {
	return []engines.Name{
		engines.BINGIMAGES,
		engines.GOOGLEIMAGES,
	}
}

func NewInfo() []engines.Name {
	return []engines.Name{
		engines.BING,
		engines.GOOGLE,
		engines.MOJEEK,
	}
}

func NewScience() []engines.Name {
	return []engines.Name{
		engines.GOOGLESCHOLAR,
	}
}

func New() Config {
	return Config{
		Server: Server{
			Port:         3030,
			FrontendUrls: []string{"http://localhost:5173"},
			Cache: Cache{
				Type: "badger",
				TTL: TTL{
					Time:        moretime.Week,
					RefreshTime: 3 * moretime.Day,
				},
				Badger: Badger{
					Persist: true,
				},
				Redis: Redis{
					Host: "localhost",
					Port: 6379,
				},
			},
			Proxy: Proxy{
				Timeouts: ProxyTimeouts{
					Dial:         3 * time.Second,
					KeepAlive:    3 * time.Second,
					TLSHandshake: 2 * time.Second,
				},
			},
		},
		Settings: NewSettings(),
		Categories: map[category.Name]Category{
			category.GENERAL: {
				Engines: NewGeneral(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.IMAGES: {
				Engines: NewImage(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1500 * time.Millisecond,
					PageTimeout: 1500 * time.Millisecond,
				},
			},
			category.INFO: {
				Engines: NewInfo(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.SCIENCE: {
				Engines: NewScience(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     3000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.NEWS: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.BLOG: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     2500 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.SURF: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     2000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.NEWNEWS: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
		},
	}
}
