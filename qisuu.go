package ebookparse

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//Qiusuu 求书中文网
type Qisuu struct {
	Name string
}

func (*Qisuu) GetBook(url string) Book {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Errorf("下载错误:%#v", err)
		return Book{}
	}
	sub := doc.Find(".showInfo").Find("p").Eq(0)
	main := doc.Find(".detail_right")
	title, _ := main.Find("h1").Html()

	sel := main.Find("li.small")

	size, _ := sel.Eq(1).Html()
	status, _ := sel.Eq(4).Html()
	author, _ := sel.Eq(5).Html()
	href, _ := sel.Eq(6).Find("a").Attr("href")
	desc, _ := sub.Html()
	size = strings.SplitAfter(size, "：")[1]
	status = strings.SplitAfter(status, "：")[1]
	author = strings.SplitAfter(author, "：")[1]
	title = strings.Replace(title, "全集", "", 1)
	href = "https://www.qisuu.la" + href

	book := Book{BookName: title, Size: size, Status: status, Author: author, Source: href, BookDecs: desc}
	return book
}

func (*Qisuu) GetDetailURLs(book Book) []Chapter {
	var osa = make([]Chapter, 0)
	doc, err := goquery.NewDocument(book.Source)
	if err != nil {
		fmt.Errorf("下载错误:%#v", err)
		os.Exit(-1)
	}
	main := doc.Find("div.pc_list").Eq(1)
	sel := main.Find("a")

	for i := range sel.Nodes {
		single := sel.Eq(i)
		href, _ := single.Attr("href")
		name, _ := single.Html()
		charID := strings.Replace(href, ".html", "", 1)
		c := Chapter{Title: name, Source: book.Source + href, BookID: book.BookID, ChapID: charID}
		osa = append(osa, c)
		//fmt.Println("find url ", i, " ", c.Source)
	}
	return osa
}

func (*Qisuu) ParseDetail(obj Chapter) Chapter {
	doc, err := goquery.NewDocument(obj.Source)
	if err != nil {
		fmt.Errorf("下载错误:%#v", err)
		//os.Exit(-1)
		return obj
	}
	main := doc.Find("div#content1")
	body, _ := main.Html()
	body = strings.Replace(body, "<br/>", "", -1)
	sel := doc.Find("div.txt_lian").Eq(1).Find("a")
	pre, _ := sel.Eq(1).Attr("href")
	next, _ := sel.Eq(3).Attr("href")
	preID := ""
	nextID := ""
	if strings.Contains(pre, "html") {
		preID = strings.Replace(pre, ".html", "", 1)
	}
	if strings.Contains(next, "html") {
		nextID = strings.Replace(next, ".html", "", 1)
	}

	obj.Content = body
	obj.PreChap = preID
	obj.NextChap = nextID

	return obj
}
