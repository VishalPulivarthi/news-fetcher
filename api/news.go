package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"news-fetcher/db"
	"news-fetcher/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// NewsResult represents a single news article in the API response
type NewsResult struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Link        string   `json:"link"`
	PubDate     string   `json:"pubDate"`
	ImageURL    string   `json:"image_url"`
	Country     []string `json:"country"`
	Category    []string `json:"category"`
	SourceID    string   `json:"source_id"`
}

// NewsAPIResponse represents the full response from the News API
type NewsAPIResponse struct {
	Status       string       `json:"status"`
	TotalResults int          `json:"totalResults"`
	Results      []NewsResult `json:"results"`
}

// FetchNews fetches news from the API and saves top articles to DB
func FetchNews(c *gin.Context) {
	var userRequest models.UserRequest
	if err := c.BindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := godotenv.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load .env file"})
		return
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API key missing"})
		return
	}

	url := fmt.Sprintf(
		"https://newsdata.io/api/1/news?apikey=%s&country=%s&language=en",
		apiKey,
		userRequest.Location,
	)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var newsResp NewsAPIResponse
	if err := json.Unmarshal(body, &newsResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}

	if newsResp.Status != "success" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API returned non-success status"})
		return
	}

	top := userRequest.TopCount
	if top > len(newsResp.Results) {
		top = len(newsResp.Results)
	}

	var articles []models.NewsArticle
	for i := 0; i < top; i++ {
		r := newsResp.Results[i]
		article := models.NewsArticle{
			Title:       r.Title,
			Description: r.Description,
			Link:        r.Link,
			PubDate:     r.PubDate,
			SourceID:    r.SourceID,
			Category:    r.Category,
			Country:     r.Country,
		}
		articles = append(articles, article)
	}

	if err := db.SaveArticles(userRequest.Location, articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save articles to DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "News fetched and stored successfully",
		"articles_saved": strconv.Itoa(len(articles)),
	})
}
