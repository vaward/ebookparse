package main

import (
	"encoding/json"
	"fmt"

	"github.com/vaward/ebookparse"
)

func main() {
	url := "http://www.quanben.co/info/3152.html"

	var parse ebookparse.BookCrawler

	parse = &ebookparse.QuanbenCO{}

	book := parse.GetBook(url)

	fmt.Println(book.Source)

	chapter := parse.GetDetailURLs(book)

	fmt.Println("get chapters count:", len(chapter))

	instance := chapter[0]

	body, _ := json.Marshal(instance)
	fmt.Println(string(body))

	fmt.Println(instance.Source)

	instance = parse.ParseDetail(instance)

	fmt.Println()
}
