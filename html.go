package html2epub

import (
	"encoding/xml"
	"fmt"
	"github.com/bmaupin/go-epub"
	"golang.org/x/net/html/charset"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func Load(url string, title string, author string, id string, person string) {
	fmt.Printf("Start - %s(%s), ", title, id)
	response, err := http.Get(url)
	if err != nil {
		return
	}

	body := response.Body

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Println(err)
		}
	}(body)
	h := html{}
	decode := xml.NewDecoder(body)
	decode.Strict = false
	decode.CharsetReader = charset.NewReaderLabel
	err = decode.Decode(&h)
	if err != nil {
		fmt.Printf("Error - id : %s, Title: %s, Error : %v\n", id, title, err)
	}

	e := epub.NewEpub(title)

	e.SetAuthor(author)

	content := h.Body.Content

	imageRegExp := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)

	subMatchSlice := imageRegExp.FindAllStringSubmatch(content, -1)
	for _, item := range subMatchSlice {
		path := item[1]
		if !strings.HasPrefix(path, "..") && !strings.HasPrefix(path, "https::") {
			i, _ := e.AddImage("https://www.aozora.gr.jp/cards/"+ person + "/files/"+path, path)
			content = strings.Replace(content, path, i, -1)
		}
	}

	content = strings.Replace(content, `<?xml version="1.0" encoding="Shift_JIS"?>`, ``, -1)
	content = strings.Replace(content, `class="notation_notes"`, `class="notation_notes" style="display:none"`, -1)
	content = strings.Replace(content, `<div id="card" style="display: block;">`, `<div id="card" style="display: none;">`, -1)
	content = strings.Replace(content, `<div id="card">`, `<div id="card" style="display: none;">`, -1)
	//content = strings.Replace(content, `src="`, `src="https://www.aozora.gr.jp/cards/` + id + "/files/", -1)

	css, err := e.AddCSS("default.css", "epub.css")
	if err != nil {
		return
	}

	_, err = e.AddSection(content, "", "", css)
	if err != nil {
		return
	}

	//err = e.Write("/Users/nhn/Documents/books/" + id + ".epub")
	//if err != nil {
	//	return
	//}

	fmt.Printf("  - end\n")
}

type html struct {
	Body body `xml:"body"`
}
type body struct {
	Content string `xml:",innerxml"`
}

