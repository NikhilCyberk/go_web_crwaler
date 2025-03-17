# Go Web Crawler

A modular web crawler implemented in Go that indexes and searches for words across web pages.

## Overview

This web crawler is designed to crawl web pages starting from a specified domain, extract text content, and create a searchable index of words found on those pages. The application is split into packages for better organization and maintainability.

## Features

- Crawls web pages within a specified domain
- Extracts text content and links from HTML pages
- Creates a searchable index of words
- Command-line interface for searching indexed words
- Concurrent crawling using goroutines and channels
- Modular design with separate packages for different responsibilities

## Project Structure

```
webcrawler/
├── main.go              # Entry point and user interface
├── crawler/
│   └── crawler.go       # Web crawling and HTML processing
├── indexer/
│   └── indexer.go       # Word indexing and search functionality
└── go.mod               # Go module definition
```

## How It Works

1. The crawler starts from a predefined URL (`AllowedDomain`)
2. It fetches the HTML content of the page
3. Extracts text and links from the HTML
4. Processes text to extract words and adds them to the index
5. Follows links to other pages within the allowed domain
6. Provides a search interface for looking up indexed words

## Communication Between Components

Components communicate using channels:
- `wordCountChan`: Passes word count information from crawler to indexer
- `visitedChan`: Tracks visited URLs to avoid duplicate processing
- `doneChan`: Signals when crawling is complete

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/webcrawler.git
cd webcrawler

# Build the project
go build
```

## Usage

```bash
# Run the web crawler
./webcrawler
```

### Menu Options

When you run the program, you'll see a menu with the following options:

1. **Search for a word** - Enter a term to search the index
2. **Exit** - Close the application

## Customization

To change the target domain for crawling, modify the constants in `crawler/crawler.go`:

```go
const AllowedDomain string = "https://your-domain.com/path/"
const BaseUrl string = "https://your-domain.com"
```

## Dependencies

This project uses only the Go standard library and does not have external dependencies.

## Concurrency Model

The application uses Go's goroutines and channels for concurrent processing:

- Main goroutine handles user interaction
- Crawler goroutine processes web pages
- Indexer goroutine manages the word index
- Channels provide thread-safe communication between components

## Future Improvements

- Add support for limiting crawl depth
- Implement persistent storage for the index
- Add more search options (phrase search, wildcard search)
- Include a web interface for searching
- Add support for multiple crawling workers
