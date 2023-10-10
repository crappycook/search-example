package main

import "gosearcher/search"

func main() {

	// Init search
	search.Init()

	// Init http server
	NewHTTPServer().Run()
	// select {}
}
