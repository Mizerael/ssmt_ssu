package parser

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly/v2"
)

func GetLinksToObject(c *colly.Collector, link string) []string {
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	var linksToObject []string
	r := regexp.MustCompile("/scp*")
	c.OnHTML("div[id=\"page-content\"]", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, h *colly.HTMLElement) {
			links := h.ChildAttrs("a", "href")
			if i > 0 && r.MatchString(links[1]) {
				linksToObject = append(linksToObject, links...)
			}
		})

	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})
	c.Visit(link)
	return linksToObject
}
