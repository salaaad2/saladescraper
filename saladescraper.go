package main

import (
	"container/list"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"fmt"
	"strings"
	"bufio"
)

func getArticles(postUrls *list.List) {
	for e := postUrls.Front(); e != nil; e = e.Next() {
		fmt.Println("INFO: writing content to : ", e.Value, ".txt")
		fmt.Println("INFO: downloading ", e.Value)
		u := e.Value.(string) // convert list content to string

		// create file
		file, err := os.Create(u[34:] + ".html")

		if err != nil {
			log.Fatal("failed to create file [", u, ".txt] skipping")
			continue
		}

		// send GET to post url
		response, err := http.Get(u)

		if err != nil {
			log.Fatal("could not download ", u, " skipping")
			defer response.Body.Close()
			continue
		}

		// convert GET response to string
		pageInBytes, err := ioutil.ReadAll(response.Body)
		postContent := string(pageInBytes)
		if err != nil {
			log.Fatal("error reading body")
			continue
		}

		// look for post content and init title
		buffer := "TMP, put title here \n"
		postStart := strings.Index(postContent, "<div class=\"body markup\">")
		postEnd   := strings.Index(postContent, "post-footer")
		if postEnd == -1 || postStart == -1 {
			log.Fatal("error : could not find post content start & end : ", postStart, "|", postEnd)
			continue
		}
		buffer += postContent[postStart:postEnd]

		// write buffer line by line to file
		paraStart := strings.Index(buffer, "<p>")
		paraEnd   := strings.Index(buffer, "</p>")
		writer := bufio.NewWriter(file)

		for paraStart != -1 && paraEnd != -1 && paraStart < paraEnd {
			writer.WriteString(buffer[paraStart:paraEnd + 4] + "\n")
			buffer = buffer[paraEnd + 4:]
			// advancein buffer
			header := strings.Index(buffer, "<h3>")
			paraStart = strings.Index(buffer, "<p>")
			if header != -1 && header < paraStart {
				paraStart = strings.Index(buffer, "<h3>")
				paraEnd   = strings.Index(buffer, "</h3>")
			} else {
				paraStart = strings.Index(buffer, "<p>")
				paraEnd   = strings.Index(buffer, "</p>")
			}
		}

		// cleanup
		writer.Flush()
		file.Close()
	}
}

func main() {
	fmt.Println("--- Welcome to this amazing web scraper. this is (for now) hardwired for substack.com ---")
	name := os.Args[1] // TODO:

	// list to hold all urls
	postUrls := list.New()

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
			postUrls.PushBack(pageContent[titleStart:titleEnd])
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
	getArticles(postUrls)
}
// <a href="https://graymirror.substack.com/p/a-clarification-on-ukraine" class="post-preview-title newsletter">
