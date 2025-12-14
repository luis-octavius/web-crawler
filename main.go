package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()

	validateArgs(len(os.Args))

	rawBaseURL := os.Args[1]

	log.Println("maxPages: ", os.Args[2])
	log.Println("maxConcurrency: ", os.Args[3])
	maxPages, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Println("couldn't parse maxPages value")
		os.Exit(1)
	}

	maxConcurrency, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Println("couldn't parse maxConcurrency value")
	}

	cfg, err := configure(rawBaseURL, maxPages, maxConcurrency)
	if err != nil {
		fmt.Printf("error - configure: %v\n", err)
		return
	}

	fmt.Printf("starting crawl of: %s...\n", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	// for normalizedURL, count := range cfg.pages {
	// 	fmt.Printf("%d - %s\n", count, normalizedURL)
	// }

	elapsed := time.Now().Sub(start)
	fmt.Printf("Execution took %v\n", elapsed)
	fmt.Printf("Execution took %f seconds\n", elapsed.Seconds())
	fmt.Printf("Execution took %d milliseconds\n", elapsed.Milliseconds())
}

func validateArgs(lenArgs int) {
	switch lenArgs {
	case 0:
		log.Println("no website provided")
		os.Exit(1)
	case 1:
		log.Println("no max pages and max concurrency provided")
		os.Exit(1)
	case 2:
		log.Println("no max concurrency provided")
		os.Exit(1)
	case 3:
		log.Println("not enough arguments provided")
		os.Exit(1)
	}
}
