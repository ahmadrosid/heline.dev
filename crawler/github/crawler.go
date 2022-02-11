package github

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/ahmadrosid/heline/solr"
	"github.com/ahmadrosid/heline/utils"
	"github.com/gocolly/colly/v2"
)

type GithuRepository struct {
	Name      string
	CountFile int
}

func NewRepository(name string) *GithuRepository {
	println("Scraping : " + name)
	return &GithuRepository{
		Name:      name,
		CountFile: 0,
	}
}

func (g *GithuRepository) CountFiles() {
	g.CountFile += 1
}

func (g *GithuRepository) Collect() error {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnHTML("html", func(h *colly.HTMLElement) {
		link := h.Request.URL.String()
		if utils.ItHas(link, "/blob/") {
			println(".")

			var file = ExtractFile(h)
			data, _ := json.Marshal([]GithubFile{file})
			payload := bytes.NewReader(data)

			err := solr.Insert(payload)
			if err != nil {
				println("Failed processing: ", link)
				return
			}

			g.CountFiles()
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		if g.isAcceptedUrl(href) {
			link := e.Request.AbsoluteURL(href)

			parsedUrl, err := url.Parse(link)
			if err != nil {
				println(err.Error())
				return
			}

			urlPaths := utils.SplitString(parsedUrl.Path, "/")
			if len(urlPaths) > 3 && urlPaths[2] == "blob" && len(urlPaths[3]) == 40 {
				return
			}

			if utils.ItHasSuffix(link, ".md") {
				link = link + "?plain=1"
			}
			c.Visit(link)
		}
	})

	return c.Visit("https://github.com/" + g.Name)
}

func isMediaFile(url string) bool {
	if strings.HasSuffix(url, ".png") {
		return true
	} else if strings.HasSuffix(url, ".jpg") {
		return true
	}

	return false
}

func (g *GithuRepository) isAcceptedUrl(url string) bool {
	accept := false
	if utils.ItHasSuffix(url, ".md?plain=1") {
		accept = false
	} else if isMediaFile(url) {
		accept = false
	} else if utils.ItHas(url, g.Name+"/tree/") {
		accept = true
	} else if utils.ItHas(url, g.Name+"/blob/") {
		accept = true
	}

	return accept
}
