package scraper

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gocolly/colly/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/useragent"
)

const searcherAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
const suggesterAccept = "application/x-suggestions+json,application/json;q=0.9,*/*;q=0.8"

func (e *EngineBase) initCollectorSearcher(ctx context.Context) {
	e.initCollector(ctx, searcherAccept)
}

func (e *EngineBase) initCollectorSuggester(ctx context.Context) {
	e.initCollector(ctx, suggesterAccept)
}

func (e *EngineBase) initCollector(ctx context.Context, acceptS string) {
	// Get a random user agent with it's Sec-CH-UA headers.
	ua := useragent.RandomUserAgentWithHeaders()

	// Initialize the collector.
	e.collector = colly.NewCollector(
		colly.StdlibContext(ctx),
		colly.Async(),
		colly.MaxDepth(1),
		colly.IgnoreRobotsTxt(),
		colly.UserAgent(ua.UserAgent),
		colly.Headers(map[string]string{
			"Accept":             acceptS,
			"Accept-Encoding":    "gzip, deflate, br, zstd",
			"Accept-Language":    "en-US,en;q=0.9",
			"Sec-Ch-Ua":          ua.SecCHUA,
			"Sec-Ch-Ua-Mobile":   ua.SecCHUAMobile,
			"Sec-Ch-Ua-Platform": ua.SecCHUAPlatform,
			"Sec-Fetch-Dest":     "document",
			"Sec-Fetch-Mode":     "navigate",
			"Sec-Fetch-Site":     "none",
		}),
	)
}

func (e *EngineBase) initCollectorOnRequest(ctx context.Context) {
	e.collector.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if IsTimeoutError(err) {
				log.Trace().
					Caller().
					Err(err).
					Str("engine", e.Name.String()).
					Msg("Context timeout error")
			} else {
				log.Error().
					Caller().
					Err(err).
					Str("engine", e.Name.String()).
					Msg("Context error")
			}
			r.Abort()
			return
		}
	})
}

func (e *EngineBase) initCollectorOnResponse() {
	e.collector.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Headers.Get("Content-Encoding"), "br") {
			reader := brotli.NewReader(bytes.NewReader(r.Body))

			body, err := io.ReadAll(reader)
			if err != nil {
				log.Error().
					Caller().
					Err(err).
					Str("engine", e.Name.String()).
					Msg("Failed to decode brotli response")
				return
			}

			r.Body = body
		} else if strings.Contains(r.Headers.Get("Content-Encoding"), "zstd") {
			reader, err := zstd.NewReader(bytes.NewReader(r.Body))
			if err != nil {
				log.Error().
					Caller().
					Err(err).
					Str("engine", e.Name.String()).
					Msg("Failed to create zstd reader")
				return
			}

			body, err := io.ReadAll(reader)
			if err != nil {
				log.Error().
					Caller().
					Err(err).
					Str("engine", e.Name.String()).
					Msg("Failed to decode zstd response")
				return
			}

			r.Body = body
		}
	})
}

func (e *EngineBase) initCollectorOnError() {
	e.collector.OnError(func(r *colly.Response, err error) {
		if IsTimeoutError(err) {
			log.Trace().
				Caller().
				// Err(err). // Timeout error produces Get "url" error with the query.
				Str("engine", e.Name.String()).
				// Str("url", urll). // Can't reliably anonymize it (because it's engine dependent).
				Msg("Request timeout error for url")
		} else {
			log.Error().
				Caller().
				Err(err).
				Str("engine", e.Name.String()).
				// Str("url", urll). // Can't reliably anonymize it (because it's engine dependent).
				Bytes("response", r.Body). // WARN: Query can be present, depending on the response from the engine.
				Msg("Request error for url")
		}
	})
}
