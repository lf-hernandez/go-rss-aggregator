package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

func startScraping(db *database.Queries,
	concurrentUnits int,
	timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v go routines every %s duration", concurrentUnits, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrentUnits),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, databaseError := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if databaseError != nil {
		log.Println("Error marking feed as fetched: ", databaseError)
	}

	rssFeed, fetchingError := urlToFeed(feed.Url)
	if fetchingError != nil {
		log.Println("Error fetching feed: ", fetchingError)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post: ", item.Title, "on feed: ", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
