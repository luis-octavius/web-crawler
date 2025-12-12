package main

import (
	"log"
	"os"
)

func main() {

	args := os.Args
	actualArgs := args[1:]
	if len(actualArgs) < 1 {
		log.Println("no website provided")
		os.Exit(1)
	}

	if len(actualArgs) > 1 {
		log.Println("too many arguments provided")
		os.Exit(1)
	}

	BASE_URL := actualArgs[0]
	log.Println("starting crawl of: ", BASE_URL)

	// html, err := getHTML(BASE_URL)
	// if err != nil {
	// 	log.Printf("error getting the html of %v: %v", BASE_URL, err)
	// 	os.Exit(1)
	// }

	pages := make(map[string]int)
	crawlPage(BASE_URL, BASE_URL, pages)

}
