package main

import (
	"fmt"
	"ssmt-ssu/parser"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("scpfoundation.net"),
	)

	linksToObject := parser.GetLinksToObject(c, "https://scpfoundation.net/scp-series")
	fmt.Printf("linksToObject: %v\n", linksToObject)

}
