package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

func matchAny(urlPath string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, urlPath)
		if matched {
			return true
		} else if err != nil {
			fmt.Println(err)
		}
	}
	return false
}

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	collector := colly.NewCollector(
		colly.Async(),
		colly.UserAgent("EMS Press Linkcheck Bot (+https://ems.press)"),
	)
	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 8})

	extensions.Referer(collector)

	// Find and visit all links on ems.press pages
	collector.OnHTML("a[href], img[src], video[src]", func(element *colly.HTMLElement) {
		url := element.Request.URL

		if url.Host == "ems.press" {
			if src := element.Attr("src"); src != "" {
				element.Request.Visit(src)
			} else {
				element.Request.Visit(element.Attr("href"))
			}
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		if request.URL.Scheme != "https" && request.URL.Scheme != "http" {
			request.Abort()
		}

		exclude := []string{
			"^\\/journals\\/.*\\/articles.*",
			"^\\/journals\\/.*\\/issues.*",
			"^\\/books\\/.*\\/.*",
		}
		include := []string{
			"^\\/journals\\/msl\\/articles.*",
			"^\\/journals\\/msl\\/issues.*",
			"^\\/books\\/esiam.*",
		}

		urlPath := request.URL.Path
		matchedExclude := matchAny(urlPath, exclude)
		matchedInclude := matchAny(urlPath, include)

		if matchedExclude && !matchedInclude {
			request.Abort()
		}
	})

	collector.OnResponse(func(response *colly.Response) {
		fmt.Printf("Checked %s\n", response.Request.URL)
	})

	collector.OnError(func(response *colly.Response, err error) {
		if response.StatusCode == 503 || response.StatusCode == 999 || response.StatusCode == 0 {
			// ignore 503 and 999 and 0 status code to avoid flaky errors
			return
		}

		exitCode = 1
		request := response.Request
		fmt.Fprintf(os.Stderr, "Error: \"%v %s\" while visiting %s (Found on: %s)\n", response.StatusCode, err, request.URL, request.Headers.Get("Referer"))
	})

	collector.Visit("https://ems.press/")
	collector.Wait()
}
