package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("--- Welcome to this amazing web scraper. this is intended for substack ---")
	name := os.Args[1]
	fmt.Println(name)

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
	titleEnd := strings.Index(pageContent, "class=\"post-preview-description\"")
	count := 0
	for titleStart != -1 {
		fmt.Println("DEBUG: start index: ", titleStart, "end index: ", titleEnd, "total length :", len(pageContent))

		if titleEnd == -1 || titleStart == -1 {
			break
		}

		pageTitle := (pageContent[titleStart + len(searchStr):titleEnd - 2])
		fmt.Println("INFO: found the following post :\n ", pageTitle)
		fmt.Println("DEBUG: start index: ", titleStart, "end index: ", titleEnd)

		pageContent = pageContent[titleEnd + len("class=\"post-preview-title newsletter\""):]

		titleStart = strings.Index(pageContent, searchStr)
		titleEnd = strings.Index(pageContent, "class=\"post-preview-title newsletter\"")

		count++
	}
	if count == 0 {
		fmt.Println("ERROR: no posts to be found")
	} else {
		fmt.Println("INFO: found [", count, "] posts.")
	}
}
// <a href="https://graymirror.substack.com/p/a-clarification-on-ukraine" class="post-preview-title newsletter">
