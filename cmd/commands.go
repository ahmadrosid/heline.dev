package cmd

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func Run(args []string) int {
	c := cli.NewCLI("heline", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"scrape github": func() (cli.Command, error) {
			return &GithubCrawlerCommand{}, nil
		},
		"scrape docset": func() (cli.Command, error) {
			return &GithubCrawlerCommand{}, nil
		},
		"server": func() (cli.Command, error) {
			return &ServerCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	return exitStatus
}
