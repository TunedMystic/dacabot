package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// PrintUsage displays the different commands for the application.
func PrintUsage(flags ...*flag.FlagSet) {
	fmt.Println("\nUsage:  dacabot COMMAND [OPTIONS]")
	fmt.Println("\nWeb Application for DACA news")
	fmt.Print("\n\n")

	for _, flagset := range flags {
		fmt.Printf("[%v]\n", flagset.Name())
		flagset.PrintDefaults()
		fmt.Println()
	}

	os.Exit(0)
}

// Parse is a helper function to parse args.
func Parse(flag *flag.FlagSet) {
	if err := flag.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}

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
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		log.Fatalf("Could not convert %v to a date.\n", dateString)
	}
	return date
}

// RunCLI parses the arguments and executes the associated action.
func RunCLI() {
	opts := CmdLineOpts{}
	todayStr := time.Now().UTC().Format("2006-01-02")

	// The 'run-server' command.
	cmdRunServer := flag.NewFlagSet("run-server", flag.ExitOnError)
	cmdRunServer.IntVar(&opts.Port, "port", 8000, "Port to run the server on")

	// The 'fetch-articles' command.
	cmdFetchArticles := flag.NewFlagSet("fetch-articles", flag.ExitOnError)
	cmdFetchArticles.StringVar(&opts.From, "from", todayStr, "Limit by PublishDate >=")
	cmdFetchArticles.StringVar(&opts.To, "to", todayStr, "Limit by PublishDate <=")

	if len(os.Args) < 2 {
		PrintUsage(cmdRunServer, cmdFetchArticles)
	}

	switch os.Args[1] {

	case "run-server":
		Parse(cmdRunServer)

		// Setup server.
		server := NewServer()
		defer server.Cleanup()

		// Setup periodic tasks.
		SetupTasks()

		// Run server.
		server.Run(opts.Port)

	case "fetch-articles":
		Parse(cmdFetchArticles)
		opts.FromDate = opts.ParseDate(opts.From)
		opts.ToDate = opts.ParseDate(opts.To)

		UpdateArticles(opts.FromDate, opts.ToDate, true)

	default:
		fmt.Printf("Unknown command %v\n", os.Args[1])
		os.Exit(1)
	}
}
