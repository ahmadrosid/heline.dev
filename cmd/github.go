package cmd

import (
	"strings"

	"github.com/ahmadrosid/heline/crawler/github"
)

type GithubCrawlerCommand struct {
}

func (c *GithubCrawlerCommand) Run(args []string) int {
	repo := ""
	if len(args) > 0 {
		repo = args[0]
	} else {
		println("Please provide github repository name!")
	}

	spider := github.NewRepository(repo)
	err := spider.Collect()
	if err != nil {
		println(err.Error())
		return 1
	}

	println("Done! Total files crawled: ", spider.CountFile)
	return 0
}

func (c *GithubCrawlerCommand) Help() string {
	helpText := `
This command will scrape github repository

Usage:
  heline scrape [option]

Option:
  repo	String github repo name ex. "ahmadrosid/gemoji".
`
	return strings.TrimSpace(helpText)
}

func (c *GithubCrawlerCommand) Synopsis() string {
	return "Scrape github repository"
}
