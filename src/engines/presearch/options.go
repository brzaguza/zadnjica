package presearch

import "github.com/hearchco/hearchco/src/engines"

var Info engines.Info = engines.Info{
	Domain:         "presearch.com",
	Name:           engines.PRESEARCH,
	URL:            "https://presearch.com/search?q=",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.GOOGLE},
}

/*
var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "div[x-data=\"searchResults(true)\"] > div.w-full > div.text-gray-300 > div > div > div",
	Link:        "div > div > a.text-results-link",
	Title:       "div > div span[x-html=\"result.title\"]",
	Description: "div[x-html*=\"result.description\"]",
}
*/

var Support engines.SupportedSettings = engines.SupportedSettings{
	SafeSearch: true,
}
