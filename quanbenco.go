package ebookparse

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
)

//QuanbenCO 全本小说网
type QuanbenCO struct {
	Name string
}

//GetBook 获取电子书概要信息
func (*QuanbenCO) GetBook(url string) Book {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Errorf("下载错误:%#v", err)
		return Book{}
	}
	trans, _ := doc.Html()
	utf8body, _ := iconv.NewReader(strings.NewReader(trans), "GB18030", "utf-8")
	doc, err = goquery.NewDocumentFromReader(utf8body)

	main := doc.Find("div#content")

	title, _ := main.Find("h1").Html()

	lis := main.Find("li")
	author, _ := main.Find("ul.novel_msg").Find("a").Html()
	desc, _ := main.Find("li#description1").Html()
	href, _ := main.Find("li.button2").Find("a").Attr("href")

	status, _ := lis.Eq(3).Find("em").Html()
	btype, _ := lis.Eq(4).Html()
	size, _ := lis.Eq(5).Html()

	desc = string([]byte(desc)[49:])
	btype = string([]byte(btype)[9:])
	size = string([]byte(size)[12:])
	href = "http://www.quanben.co" + href

	book := Book{BookName: title, Size: size, Status: status, Author: author, Source: href, BookType: btype, BookDecs: desc}
	return book
}

//GetDetailURLs 获取章节信息列表
func (*QuanbenCO) GetDetailURLs(book Book) []Chapter {
	var osa = make([]Chapter, 0)
	doc, err := goquery.NewDocument(book.Source)
	if err != nil {
		fmt.Errorf("下载错误:%#v", err)
		os.Exit(-1)
	}
	trans, _ := doc.Html()
	utf8body, _ := iconv.NewReader(strings.NewReader(trans), "GBK", "utf-8")
	doc, err = goquery.NewDocumentFromReader(utf8body)

	main := doc.Find("div.novel_volume")
	sel := main.Find("a")

	for i := range sel.Nodes {
		single := sel.Eq(i)
		href, _ := single.Attr("href")
		name, _ := single.Html()
		chapID := strings.Replace(href, ".html", "", 1)

		c := Chapter{Title: name, Source: strings.Replace(book.Source, "index", chapID, 1), BookID: book.BookID, ChapID: chapID}
		osa = append(osa, c)
		//fmt.Println("find url ", i, " ", c.Source)
	}
	return osa
}

//ParseDetail 解析章节详情
func (*QuanbenCO) ParseDetail(obj Chapter) Chapter {
	fmt.Println(obj.Source)
	doc, err := goquery.NewDocument(obj.Source)
	if err != nil {
		fmt.Errorf("下载错误:%#v", err)
		//os.Exit(-1)
		return obj
	}
	trans, _ := doc.Html()
	trans = strings.Replace(trans, "›", "", -1)
	result, _ := iconv.ConvertString(trans, "gb18030", "utf-8")
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(result))

	main := doc.Find("div.novel_content")
	body, _ := main.Html()
	body = strings.Replace(body, "<br/>", "", -1)
	sel := doc.Find("div.novel_bottom_wap").Find("a")
	pre, _ := sel.Eq(0).Attr("href")

	next, _ := sel.Eq(2).Attr("href")
	preID := ""
	nextID := ""
	ps := strings.Split(pre, "/")
	pre = ps[len(ps)-1]
	ps = strings.Split(next, "/")
	next = ps[len(ps)-1]
	if !strings.EqualFold(pre, "index.html") {
		preID = strings.Replace(pre, ".html", "", 1)
	}
	if !strings.EqualFold(next, "index.html") {
		nextID = strings.Replace(next, ".html", "", 1)
	}
	fmt.Println(body)
	body = strings.Split(body, "</div>")[1]
	body = strings.Replace(body, "聽", " ", -1)

	obj.Content = body
	obj.PreChap = preID
	obj.NextChap = nextID

	return obj
}
