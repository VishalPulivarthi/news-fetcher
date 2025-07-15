package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"news-fetcher/models"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes the database connection and creates the news table if it does not exist.
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "news.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Verify DB connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS news (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		location TEXT,
		title TEXT,
		description TEXT,
		link TEXT,
		pub_date TEXT,
		source_id TEXT,
		category TEXT,
		country TEXT
	);
	`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Database initialized successfully.")
}

// SaveArticles inserts fetched news articles into the database.
func SaveArticles(location string, articles []models.NewsArticle) error {
	// Prepare the insert statement once for efficiency
	stmt, err := DB.Prepare(`
		INSERT INTO news (location, title, description, link, pub_date, source_id, category, country)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	for _, article := range articles {
		// Convert slice fields to comma-separated strings if needed
		category := strings.Join(article.Category, ",")
		country := strings.Join(article.Country, ",")

		_, err := stmt.Exec(
			location, article.Title, article.Description, article.Link,
			article.PubDate, article.SourceID, category, country,
		)
		if err != nil {
			return fmt.Errorf("failed to insert article: %w", err)
		}
	}
	fmt.Printf("✅ %d articles saved to database.\n", len(articles))
	return nil
}
