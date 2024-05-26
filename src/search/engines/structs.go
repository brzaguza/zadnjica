package engines

import (
	"fmt"

	"github.com/hearchco/hearchco/src/search/category"
)

type SupportedSettings struct {
	Locale                  bool
	SafeSearch              bool
	Mobile                  bool
	RequestedResultsPerPage bool
}

type Info struct {
	Name           Name
	ResultsPerPage int
	Domain         string
	URL            string
}

type DOMPaths struct {
	ResultsContainer string
	Result           string // div
	Link             string // a href
	Title            string // heading
	Description      string // paragraph
}

type Pages struct {
	Start int
	Max   int
}

type Options struct {
	VisitPages bool
	SafeSearch bool
	Mobile     bool
	Pages      Pages
	UserAgent  string
	Locale     string //format: en_US
	Category   category.Name
}

func ValidateLocale(locale string) error {
	if locale == "" {
		return nil
	}

	if len(locale) != 5 {
		return fmt.Errorf("engines.validateLocale(): isn't 5 characters long")
	}
	if !(('a' <= locale[0] && locale[0] <= 'z') && ('a' <= locale[1] && locale[1] <= 'z')) {
		return fmt.Errorf("engines.validateLocale(): first two characters must be lowercase ASCII letters")
	}
	if !(('A' <= locale[3] && locale[3] <= 'Z') && ('A' <= locale[4] && locale[4] <= 'Z')) {
		return fmt.Errorf("engines.validateLocale(): last two characters must be uppercase ASCII letters")
	}
	if locale[2] != '_' {
		return fmt.Errorf("engines.validateLocale(): third character must be underscore (_)")
	}

	return nil
}
