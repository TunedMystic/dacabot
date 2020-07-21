package main

import "github.com/tunedmystic/dacabot/app"

func main() {
	// Setup server.
	server := app.NewServer()
	defer server.Cleanup()

	// Setup periodic tasks.
	app.SetupTasks()

	// Run server.
	server.Run()
}
