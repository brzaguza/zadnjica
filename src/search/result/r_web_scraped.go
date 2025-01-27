package result

import (
	"github.com/hearchco/agent/src/utils/moreurls"
	"github.com/rs/zerolog/log"
)

type WebScraped struct {
	url         string
	title       string
	description string
	rank        RankScraped
}

func (r WebScraped) Key() string {
	return r.URL()
}

func (r WebScraped) URL() string {
	if r.url == "" {
		log.Panic().Msg("url is empty")
		// ^PANIC - Assert because the url should never be empty.
	}

	return r.url
}

func (r WebScraped) Title() string {
	if r.title == "" {
		log.Panic().Msg("title is empty")
		// ^PANIC - Assert because the title should never be empty.
	}

	return r.title
}

func (r WebScraped) Description() string {
	return r.description
}

func (r WebScraped) Rank() RankScraped {
	return r.rank
}

func (r WebScraped) Convert(erCap int) Result {
	engineRanks := make([]Rank, 0, erCap)
	engineRanks = append(engineRanks, r.Rank().Convert())
	return &Web{
		webJSON{
			URL:         r.URL(),
			FQDN:        moreurls.FQDN(r.URL()),
			Title:       r.Title(),
			Description: r.Description(),
			EngineRanks: engineRanks,
		},
	}
}
