package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"
)

func startScraping(db *database.Queries,
	concurrentUnits int,
	timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines ('lightweight threads') every %s", concurrentUnits, timeBetweenRequest)
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
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		parsedPubAtDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
		}

		_, createPostDbError := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: parsedPubAtDate,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if createPostDbError != nil {
			if strings.Contains(createPostDbError.Error(), "duplicate key") {
				continue
			}
			fmt.Println("error creating post in database: ", createPostDbError)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
