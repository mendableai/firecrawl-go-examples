// change the api key to your own
// go run example.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mendableai/firecrawl-go"
)

func ptr[T any](v T) *T {
	return &v
}

func main() {
	app, err := firecrawl.NewFirecrawlApp("fc-YOUR_API_KEY", "https://api.firecrawl.dev")
	if err != nil {
		log.Fatalf("Failed to create FirecrawlApp: %v", err)
	}

	// Scrape a website
	scrapeResult, err := app.ScrapeURL("firecrawl.dev", nil)
	if err != nil {
		log.Fatalf("Failed to scrape URL: %v", err)
	}
	fmt.Println(scrapeResult.Markdown)

	// Crawl a website
	idempotencyKey := uuid.New().String() // optional idempotency key
	crawlParams := &firecrawl.CrawlParams{
		ExcludePaths: []string{"blog/*"},
		MaxDepth:     ptr(2),
	}
	crawlResult, err := app.CrawlURL("mendable.ai", crawlParams, &idempotencyKey)
	if err != nil {
		log.Fatalf("Failed to crawl URL: %v", err)
	}
	jsonCrawlResult, err := json.MarshalIndent(crawlResult, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal crawl result: %v", err)
	}
	fmt.Println(string(jsonCrawlResult))

	asyncCrawlParams := &firecrawl.CrawlParams{
		ExcludePaths: []string{"blog/*"},
		ScrapeOptions: firecrawl.ScrapeParams{
			Formats: []string{"markdown", "html", "rawHtml", "screenshot", "links"},
		},
		MaxDepth: ptr(2),
	}
	asyncCrawlResponse, err := app.AsyncCrawlURL("mendable.ai", asyncCrawlParams, nil)
	if err != nil {
		log.Fatalf("Failed to async crawl URL: %v", err)
	}

	const maxChecks = 15
	checks := 0

	for {
		if checks >= maxChecks {
			break
		}

		time.Sleep(2 * time.Second) // wait for 2 seconds

		response, err := app.CheckCrawlStatus(asyncCrawlResponse.ID)
		if err != nil {
			log.Fatalf("Failed to check crawl status: %v", err)
		}

		if response.Status == "completed" {
			break
		}

		checks++
	}

	// Final check after loop or if completed
	completedCrawlResponse, err := app.CheckCrawlStatus(asyncCrawlResponse.ID)
	if err != nil {
		log.Fatalf("Failed to check crawl status: %v", err)
	}

	jsonCompletedCrawlResponse, err := json.MarshalIndent(completedCrawlResponse, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal async crawl result: %v", err)
	}

	fmt.Println(string(jsonCompletedCrawlResponse))

	// LLM Extraction using JSON schema
	// jsonSchema := map[string]any{
	// 	"type": "object",
	// 	"properties": map[string]any{
	// 		"top": map[string]any{
	// 			"type": "array",
	// 			"items": map[string]any{
	// 				"type": "object",
	// 				"properties": map[string]any{
	// 					"title":       map[string]string{"type": "string"},
	// 					"points":      map[string]string{"type": "number"},
	// 					"by":          map[string]string{"type": "string"},
	// 					"commentsURL": map[string]string{"type": "string"},
	// 				},
	// 				"required": []string{"title", "points", "by", "commentsURL"},
	// 			},
	// 			"minItems":    5,
	// 			"maxItems":    5,
	// 			"description": "Top 5 stories on Hacker News",
	// 		},
	// 	},
	// 	"required": []string{"top"},
	// }

	// llmExtractionParams := &firecrawl.ScrapeParams{
	// 	ExtractorOptions: firecrawl.ExtractorOptions{
	// 		ExtractionSchema: jsonSchema,
	// 		Mode:             "llm-extraction",
	// 	},
	// 	PageOptions: firecrawl.PageOptions{
	// 		OnlyMainContent: prt(true),
	// 	},
	// }

	// llmExtractionResult, err := app.ScrapeURL("https://news.ycombinator.com", llmExtractionParams)
	// if err != nil {
	// 	log.Fatalf("Failed to perform LLM extraction: %v", err)
	// }

	// // Pretty print the LLM extraction result
	// jsonResult, err := json.MarshalIndent(llmExtractionResult.LLMExtraction, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Failed to marshal LLM extraction result: %v", err)
	// }
	// fmt.Println(string(jsonResult))

	mapResult, err := app.MapURL("https://firecrawl.dev", &firecrawl.MapParams{
		Search: ptr("blog"),
	})
	if err != nil {
		log.Fatalf("Failed to map URL: %v", err)
	}
	jsonMapResult, err := json.MarshalIndent(mapResult, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal map result: %v", err)
	}
	fmt.Println(string(jsonMapResult))
}
