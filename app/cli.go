package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/integrii/flaggy"
)

// CmdLineOpts stores the options that are parsed.
type CmdLineOpts struct {
	Port     int
	From     string
	To       string
	FromDate time.Time
	ToDate   time.Time
}

// ParseDate accepts a date string and returns a time.Time value.
func (c CmdLineOpts) ParseDate(dateString string) time.Time {
	return MustParseDate(dateString)
}

// RunCLI parses the given args and executes the appropriate action.
func RunCLI() {
	todayStr := time.Now().UTC().Format("2006-01-02")

	opts := CmdLineOpts{
		Port: 8000,
		From: todayStr,
		To:   todayStr,
	}

	flaggy.SetName("dacabot")
	flaggy.SetDescription("A news aggregator for DACA-related news")

	// The 'run-server' subcommand.
	cmdRunServer := flaggy.NewSubcommand("run-server")
	cmdRunServer.Description = "Start the web application"
	cmdRunServer.Int(&opts.Port, "p", "port", "Port to run the server on")
	flaggy.AttachSubcommand(cmdRunServer, 1)

	// The 'fetch-articles' subcommand.
	cmdFetchArticles := flaggy.NewSubcommand("fetch-articles")
	cmdFetchArticles.Description = "Fetch articles from news sources"
	cmdFetchArticles.String(&opts.From, "f", "from", "Limit by PublishDate >=")
	cmdFetchArticles.String(&opts.To, "t", "to", "Limit by PublishDate <=")
	flaggy.AttachSubcommand(cmdFetchArticles, 1)

	flaggy.Parse()

	if len(os.Args) < 2 {
		flaggy.ShowHelp("")
		return
	}

	if cmdRunServer.Used {
		// Setup server.
		server := NewServer()
		defer server.Cleanup()

		// Setup periodic tasks.
		SetupTasks()

		// Run server.
		// server.Run(opts.Port)

		addr := fmt.Sprintf("0.0.0.0:%v", opts.Port)
		fmt.Printf("[run] starting Server on %v...\n", addr)
		log.Fatal(http.ListenAndServe(addr, server.Router))
	}

	if cmdFetchArticles.Used {
		opts.FromDate = opts.ParseDate(opts.From)
		opts.ToDate = opts.ParseDate(opts.To)

		// Fetch articles.
		UpdateArticles(opts.FromDate, opts.ToDate, true)
	}
}
