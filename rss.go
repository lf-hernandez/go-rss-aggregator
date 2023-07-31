package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	httpResponse, httpError := httpClient.Get(url)
	if httpError != nil {
		return RSSFeed{}, httpError
	}
	defer httpResponse.Body.Close()

	data, dataError := io.ReadAll(httpResponse.Body)
	if dataError != nil {
		return RSSFeed{}, dataError
	}

	rssFeed := RSSFeed{}

	xmlParsingError := xml.Unmarshal(data, &rssFeed)
	if xmlParsingError != nil {
		return RSSFeed{}, xmlParsingError
	}

	return rssFeed, nil
}
