package ebookparse

//Book 电子书库
type Book struct {
	BookID   string `json:"objectId,omitempty"`
	BookName string `json:"name,omitempty"`
	Author   string `json:"author,omitempty"`
	BookDecs string `json:"desc,omitempty"`
	BookType string `json:"type,omitempty"`
	Status   string `json:"status,omitempty"`
	Size     string `json:"size,omitempty"`
	Source   string `json:"source,omitempty"`
}

//Chapter 目录
type Chapter struct {
	BookID   string `json:"bookId,omitempty"`
	ChapID   string `json:"chapter,omitempty"`
	Title    string `json:"title,omitempty"`
	Source   string `json:"source,omitempty"`
	Content  string `json:"content,omitempty"`
	PreChap  string `json:"pre,omitempty"`
	NextChap string `json:"next,omitempty"`
}

//BookCrawler 抓取接口
type BookCrawler interface {
	GetBook(url string) Book
	GetDetailURLs(book Book) []Chapter
	ParseDetail(chapter Chapter) Chapter
}
