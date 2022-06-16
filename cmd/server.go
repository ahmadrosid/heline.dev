package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	ghttp "github.com/ahmadrosid/heline/http"
)

type ServerCommand struct {
}

func (c *ServerCommand) Run(args []string) int {
	handler := ghttp.Handler()

	port := "80"
	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
		Addr:              ":" + port,
	}

	command := ""
	if len(args) > 0 {
		command = args[0]
	}

	switch command {
	case "start":
		fmt.Printf("🚀 Starting server on http://localhost:%s\n", port)
		err := server.ListenAndServe()
		if err != nil {
			println("❌ Server already started!")
			return 1
		}
	case "stop":
		if runtime.GOOS == "windows" {
			command := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess -Force", port)
			exec_cmd(exec.Command("Stop-Process", "-Id", command))
		} else {
			command := fmt.Sprintf("lsof -i tcp:%s | grep LISTEN | awk '{print $2}' | xargs kill -9", port)
			exec_cmd(exec.Command("bash", "-c", command))
		}

		return 0
	default:
		println(c.Help())
	}

	return 0
}

func (c *ServerCommand) Help() string {
	helpText := `
This command will manage backend rest api server.

Usage: 
  heline server [option]

Options:
  start     Start server
  stop      Stop server
  -help -h  Show help
`
	return strings.TrimSpace(helpText)
}

func (c *ServerCommand) Synopsis() string {
	return "Run backend server"
}

func exec_cmd(cmd *exec.Cmd) {
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if err != nil {
			os.Stderr.WriteString("Server not started yet.\n")
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("Server successfully stopped. (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}
