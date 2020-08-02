package main

import (
	"log"

	"github.com/tunedmystic/dacabot/app"
)

func init() {
	log.SetFlags(0)
}

func main() {

	// Execute the CLI to do things like
	// run the server or fetch articles.
	app.RunCLI()

}
