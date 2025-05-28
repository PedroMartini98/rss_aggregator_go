package scrapper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PedroMartini98/rss_aggregator_go/internal/database"
	"github.com/google/uuid"
)

func StartScrapping(db *database.Queries, concurrencyLimit int, timeBtwnRequest time.Duration) {
	log.Printf("Scrapping on %v goroutines every %s durations", concurrencyLimit, timeBtwnRequest)
	ticker := time.NewTicker(timeBtwnRequest)
	//intialize the for loop with empty ; ; so it starts imidiatly
	for ; ; <-ticker.C {

		feeds, err := db.GetFeedsToFetch(context.Background(), int32(concurrencyLimit))

		if err != nil {
			log.Printf("Error fetching feeds:%v", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapFeed(db, feed, wg)
		}
	}

}

func scrapFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {

	defer wg.Done()

	_, err := db.MarkFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("failed to mark feed( %s ) as fetched: %v ", feed.Name, err)
		return
	}

	rssFeed, err := UrlIntoFeed(feed.Url)
	if err != nil {
		log.Printf("failed to fecth feed( %s ) : %v ", feed.Name, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("failed to parse date %v with err: %v", item.PubDate, err)
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			PublishedAt: parsedTime,
			Description: description,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Print("failed to create post:", err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
