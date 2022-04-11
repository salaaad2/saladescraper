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

func getArticle(postUrl string, path string, done chan int) {
	fmt.Println("INFO: writing content to : ", postUrl, ".txt")
	fmt.Println("INFO: downloading ", postUrl)

	// create file
	file, _ := os.Create("./" + path + "/" + postUrl[34:] + ".html")
	// send GET to post url
	response, err := http.Get(postUrl)
	if err != nil {
		log.Fatal("could not download ", postUrl, " skipping")
		defer response.Body.Close()
		return
	}

	// convert GET response to string
	pageInBytes, _ := ioutil.ReadAll(response.Body)
	postContent := string(pageInBytes)

	// look for post content and init title
	buffer := "TMP, put title here \n"
	postStart := strings.Index(postContent, "<div class=\"single-post-container\">")
	postEnd   := strings.Index(postContent, "post-footer")
	if postEnd == -1 || postStart == -1 {
		log.Fatal("error : could not find post content start & end : ", postStart, "|", postEnd)
		return
	}
	buffer += postContent[postStart:postEnd]

	// write buffer line by line to file
	paraStart := strings.Index(buffer, "<p>")
	paraEnd   := strings.Index(buffer, "</p>")
	writer := bufio.NewWriter(file)
	writer.WriteString("<title> " + strings.ToUpper(postUrl[34:]) + " </title>\n")
	writer.WriteString("<h1> " + strings.ToUpper(postUrl[34:]) + " </h1>\n")

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
	done <- 1
}

func main() {
	fmt.Println("--- Welcome to this amazing web scraper. this is (for now) hardwired for substack.com ---")
	var target string
	if len(os.Args) <= 1 {
		log.Fatal("missing argument: substack url")
		return
	}
	target = os.Args[1]

	// list to hold all urls
	postUrls := list.New()

	// create url from parameter
	one := "https://"
	two := one + target
	three := ".substack.com"
	full := two + three

	searchStr := full + "/p"
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
		log.Fatal("ERROR: no posts to be found")
		return
	} else {
		fmt.Printf("found [%d] posts.", count)
	}
	// create directory
	dir_err := os.MkdirAll(target, 0755)
	if dir_err != nil {
		log.Fatal("failed to create directory", target)
		os.Exit(1)
	}
	done := make(chan int)
	for e := postUrls.Front(); e != nil; e = e.Next() {
		fmt.Println("getting article")
		go getArticle(e.Value.(string), target, done)
	}

	for {
		select {
		case <-done:
			count--
			if count == 0 {
				fmt.Println("finished downloading substack posts")
				return
			}
		}
	}
}
