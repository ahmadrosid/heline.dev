package cmd

import (
	"strings"

	"github.com/ahmadrosid/heline/crawler/docset"
)

type DocsetCrawlerCommand struct {
}

func (c *DocsetCrawlerCommand) Run(args []string) int {
	// repo := ""
	// if len(args) > 0 {
	// 	repo = args[0]
	// } else {
	// 	println("Please provide docset name!")
	// }

	spider := docset.NewDocset()
	err := spider.Collect()
	if err != nil {
		println(err.Error())
		return 1
	}

	println("Done! Total files crawled: ", spider.CountFile)
	return 0
}

func (c *DocsetCrawlerCommand) Help() string {
	helpText := `
This command will scrape github repository

Usage:
  heline scrape [option]

Option:
  repo	String github repo name ex. "ahmadrosid/gemoji".
`
	return strings.TrimSpace(helpText)
}

func (c *DocsetCrawlerCommand) Synopsis() string {
	return "Scrape github repository"
}
