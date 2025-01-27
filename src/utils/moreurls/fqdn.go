package moreurls

import (
	"net/url"

	"github.com/rs/zerolog/log"
)

// Returns the fully qualified domain name of the URL.
func FQDN(urll string) string {
	// Check if the url is empty.
	if urll == "" {
		log.Panic().
			Str("url", urll).
			Msg("URL is empty")
	}

	// Parse the URL.
	u, err := url.Parse(urll)
	if err != nil {
		log.Panic().
			Err(err).
			Str("url", urll).
			Msg("Failed to parse the URL")
		// ^PANIC - Assert correct URL.
	}

	// Check if the hostname is empty.
	h := u.Hostname()
	if h == "" {
		log.Panic().
			Str("url", urll).
			Msg("Hostname is empty")
		// ^PANIC - Assert non-empty URL.
	}

	return h
}
