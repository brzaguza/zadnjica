package search

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/utils/anonymize"
)

// Searches for suggestions using the provided category config.
func Suggestions(query string, opts options.Options, catConf category.Category) ([]result.Suggestion, error) {
	// Capture start time.
	startTime := time.Now()

	if err := validateParams(query, opts); err != nil {
		return nil, err
	}

	log.Debug().
		Str("query", anonymize.String(query)).
		Str("locale", opts.Locale.String()).
		Msg("Suggesting")

	// Create contexts with timeout for HardTimeout and PreferredTimeout.
	ctxHardTimeout, cancelHardTimeoutFunc := context.WithTimeout(context.Background(), catConf.Timings.HardTimeout)
	defer cancelHardTimeoutFunc()
	ctxPreferredTimeout, cancelPreferredTimeoutFunc := context.WithTimeout(context.Background(), catConf.Timings.PreferredTimeout)
	defer cancelPreferredTimeoutFunc()

	// Create a context that cancels when both HardTimeout and PreferredTimeout are done.
	suggestCtx, cancelSuggest := context.WithCancel(context.Background())
	defer cancelSuggest()
	go func() {
		<-ctxHardTimeout.Done()
		<-ctxPreferredTimeout.Done()
		cancelSuggest()
	}()

	// Initialize each engine.
	suggesters := initializeSuggesters(suggestCtx, catConf.Engines)

	// Create a map for the suggestions with RWMutex.
	concMap := result.NewSuggestionMap(len(catConf.Engines))

	// Create a sync.Once wrapper for each suggester.Suggest() to ensure that the engine is only run once.
	onceWrapMap := initOnceWrapper(catConf.Engines)

	// Run all required engines. WaitGroup should be awaited unless the hard timeout is reached.
	var wgRequiredEngines sync.WaitGroup
	runRequiredSuggesters(catConf.RequiredEngines, suggesters, &wgRequiredEngines, &concMap, query, opts, onceWrapMap)

	// Run all required by origin engines. Cond should be awaited unless the hard timeout is reached.
	var wgRequiredByOriginEngines sync.WaitGroup
	runRequiredByOriginSuggesters(catConf.RequiredByOriginEngines, suggesters, &wgRequiredByOriginEngines, &concMap, catConf.Engines, query, opts, onceWrapMap)

	// Run all preferred engines. WaitGroup should be awaited unless the preferred timeout is reached.
	var wgPreferredEngines sync.WaitGroup
	runPreferredSuggesters(catConf.PreferredEngines, suggesters, &wgPreferredEngines, &concMap, query, opts, onceWrapMap)

	// Run all preferred by origin engines. Cond should be awaited unless the preferred timeout is reached.
	var wgPreferredByOriginEngines sync.WaitGroup
	runPreferredByOriginSuggesters(catConf.PreferredByOriginEngines, suggesters, &wgPreferredByOriginEngines, &concMap, catConf.Engines, query, opts, onceWrapMap)

	// Cancel the hard timeout after all required engines have finished and all required by origin engines have finished.
	go cancelHardTimeout(startTime, cancelHardTimeoutFunc, query, &wgRequiredEngines, catConf.RequiredEngines, &wgRequiredByOriginEngines, catConf.RequiredByOriginEngines)

	// Cancel the preferred timeout after all preferred engines have finished and all preferred by origin engines have finished.
	go cancelPreferredTimeout(startTime, cancelPreferredTimeoutFunc, query, &wgPreferredEngines, catConf.PreferredEngines, &wgPreferredByOriginEngines, catConf.PreferredByOriginEngines)

	// Wait for both hard timeout and preferred timeout to finish.
	<-suggestCtx.Done()

	// Extract the suggestions and responders from the map.
	suggestions, responders := concMap.ExtractWithResponders()

	log.Debug().
		Int("suggestions", len(suggestions)).
		Str("query", anonymize.String(query)).
		Str("responders", fmt.Sprintf("%v", responders)).
		Dur("duration", time.Since(startTime)).
		Msg("Scraping finished")

	// Return the suggestions.
	return suggestions, nil
}
