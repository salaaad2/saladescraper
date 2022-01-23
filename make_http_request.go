package main

import (
	"container/list"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("--- Welcome to this amazing web scraper. this is intended for substack ---")
	name := os.Args[1] // TODO:


	// create url from parameter
	one := "https://"
	two := one + name
	three := ".substack.com"
	full := two + three

	searchStr := full + "/p/"

	fmt.Println("INFO: Trying to get all articles from : ", full)

	response, err := http.Get(full)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// get all content and convert it into a string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)
	if err != nil {
		log.Fatal(err)
	}

	// look for posts
	titleStart := strings.Index(pageContent, searchStr)
	titleEnd := strings.Index(pageContent, "/comments")
	count := 0 // how many posts

	for titleEnd != -1 && titleStart != -1 { // go through whole text

		if (titleEnd - titleStart) > 100 {
			pageContent = pageContent[titleStart + len(searchStr):]
			titleStart = strings.Index(pageContent, searchStr)
			titleEnd = strings.Index(pageContent, "/comments")
		} else {
			fmt.Println(pageContent[titleStart:titleEnd])
			pageContent = pageContent[titleEnd + len(searchStr + "/comments"):]
			titleStart = strings.Index(pageContent, searchStr)
			titleEnd = strings.Index(pageContent, "/comments")
			count++
		}
	}
	if count == 0 {
		fmt.Println("ERROR: no posts to be found")
	} else {
		fmt.Println("INFO: found [", count, "] posts.")
	}
}
// <a href="https://graymirror.substack.com/p/a-clarification-on-ukraine" class="post-preview-title newsletter">
