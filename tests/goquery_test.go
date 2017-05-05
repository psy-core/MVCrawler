package tests

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"log"
	"testing"
)

func ExampleScrape() {
	doc, err := goquery.NewDocument("http://mv.yinyuetai.com/all.html#sid=5%3B12&tid=33%3B73&a=&p=&c=sh")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".mv-list > a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band, exist := s.Find("a").Attr("href")
		fmt.Printf("exist:%v\n", exist)
		fmt.Printf("Review %d: %s \n", i, band)
	})
}

func TestGoquery(t *testing.T) {
	ExampleScrape()
}
