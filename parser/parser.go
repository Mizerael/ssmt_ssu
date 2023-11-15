package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ScpObject struct {
	Name     string `json:"name"`
	Class    string `json:"class"`
	Contains string `json:"contains"`
}

const LinkToSite = "scpfoundation.net"
const DIR = "data/"

var rLinks = regexp.MustCompile("/scp.*")
var navBar = regexp.MustCompile("«.*|.*»")

func getLinksToObject(c *colly.Collector, link string) ([]string, error) {
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	var linksToObject []string

	c.OnHTML("div[id=\"page-content\"]", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, h *colly.HTMLElement) {
			links := h.ChildAttrs("a", "href")
			if i > 0 && rLinks.MatchString(links[1]) {
				linksToObject = append(linksToObject, links...)
			}
		})

	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	err := c.Visit("https://" + LinkToSite + link)

	return linksToObject, err
}

func getObject(c *colly.Collector, link string) (ScpObject, error) {
	var scp ScpObject
	var describe string
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnHTML("div[id=\"page-content\"]", func(e *colly.HTMLElement) {

		e.ForEach("p", func(i int, h *colly.HTMLElement) {
			switch i {
			case 0:
				splitText := strings.Split(h.Text, ":")
				if len(splitText) > 1 {
					scp.Name = splitText[1]
				} else {
					scp.Name = strings.Split(link, "\\")[0]
				}
			case 1:
				splitText := strings.Split(h.Text, ":")
				if len(splitText) > 1 {
					scp.Class = splitText[1]
				} else {
					scp.Class = "Данные удалены"
				}
			default:
				if !navBar.MatchString(h.Text) {
					describe += h.Text
				}

			}
		})
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	err := c.Visit("https://" + LinkToSite + link)

	scp.Contains = describe
	return scp, err
}

func getListOfObject(links []string) ([]ScpObject, error) {

	var scpList []ScpObject

	for _, v := range links {
		coll := colly.NewCollector(
			colly.AllowedDomains(LinkToSite),
		)
		scp, err := getObject(coll, v)
		if err != nil {
			return scpList, err
		}
		scpList = append(scpList, scp)
	}

	return scpList, nil
}

func scpToJson(scpList []ScpObject, name string) error {
	content, err := json.Marshal(scpList)
	if err != nil {
		return err
	}
	println(len(scpList))
	link := DIR + name + ".json"
	os.WriteFile(link, content, 0644)
	return nil
}

func ParsePage(link string) error {
	c := colly.NewCollector(
		colly.AllowedDomains(LinkToSite),
	)

	linksToObjects, err := getLinksToObject(c, link)
	if err != nil {
		return err
	}

	scpList, err := getListOfObject(linksToObjects)
	if err != nil {
		return err
	}

	err = scpToJson(scpList, "scp-0")
	if err != nil {
		return err
	}
	return nil
}
