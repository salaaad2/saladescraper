package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"fmt"
)

func main() {
	fmt.Println("Welcome to this amazing web scraper. this is intended for substack")
	name := os.Args[1]
	fmt.Println(name)

	one := "https://"
	two := one + name
	three := ".substack.com"
	full := two + three

	fmt.Println("Trying to get all articles from : ", full)

	response, err := http.Get(full)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	n, err := io.Copy(os.Stdout, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("bytes : ", n)
}
