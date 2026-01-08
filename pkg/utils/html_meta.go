package utils

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MetaInfo struct {
	Title       string
	Description string
	Image       string
}

func GetMetaInfo(url string) (*MetaInfo, error) {
	// Fetch the HTML
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	meta := &MetaInfo{}

	// Try OpenGraph tags first, then fall back to standard meta tags
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, _ := s.Attr("property")
		name, _ := s.Attr("name")
		content, _ := s.Attr("content")

		switch {
		case property == "og:title" || name == "og:title" || name == "title":
			meta.Title = content
		case property == "og:description" || name == "og:description":
			meta.Description = content
		case property == "og:image" || name == "og:image" || name == "twitter:image":
			meta.Image = content
		case name == "description" && meta.Description == "":
			meta.Description = content
		}
	})

	// Fallback to <title> tag if no og:title found
	if meta.Title == "" {
		meta.Title = strings.TrimSpace(doc.Find("title").First().Text())
	}

	return meta, nil
}
