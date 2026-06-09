package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const newsAPIBaseURL = "https://newsapi.org/v2"

// Client is the newsapi.org HTTP client.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func newClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: newsAPIBaseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// --- Response types ---

type ArticlesResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
	Message      string    `json:"message,omitempty"`
	Code         string    `json:"code,omitempty"`
}

type Article struct {
	Source      ArticleSource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	URLToImage  string        `json:"urlToImage"`
	PublishedAt string        `json:"publishedAt"`
	Content     string        `json:"content"`
}

type ArticleSource struct {
	ID   *string `json:"id"`
	Name string  `json:"name"`
}

type SourcesResponse struct {
	Status  string   `json:"status"`
	Sources []Source `json:"sources"`
	Message string   `json:"message,omitempty"`
	Code    string   `json:"code,omitempty"`
}

type Source struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	Country     string `json:"country"`
}

// --- Parameter types ---

type TopHeadlinesParams struct {
	Country  string
	Category string
	Sources  string
	Q        string
	QInTitle string
	Language string
	PageSize int
	Page     int
}

type EverythingParams struct {
	Q              string
	QInTitle       string
	SearchIn       string
	Sources        string
	Domains        string
	ExcludeDomains string
	From           string
	To             string
	Language       string
	SortBy         string
	PageSize       int
	Page           int
}

type SourcesParams struct {
	Category string
	Language string
	Country  string
}

// --- Client methods ---

func (c *Client) get(endpoint string, params url.Values) ([]byte, error) {
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("User-Agent", "newsapi-mcp/1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) GetTopHeadlines(p TopHeadlinesParams) (*ArticlesResponse, error) {
	params := url.Values{}
	setString(params, "country", p.Country)
	setString(params, "category", p.Category)
	setString(params, "sources", p.Sources)
	setString(params, "q", p.Q)
	setString(params, "qintitle", p.QInTitle)
	setString(params, "language", p.Language)
	setInt(params, "pageSize", p.PageSize)
	setInt(params, "page", p.Page)

	body, err := c.get("/top-headlines", params)
	if err != nil {
		return nil, err
	}

	var result ArticlesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if result.Status == "error" {
		return nil, apiError(result.Code, result.Message)
	}
	return &result, nil
}

func (c *Client) SearchEverything(p EverythingParams) (*ArticlesResponse, error) {
	params := url.Values{}
	setString(params, "q", p.Q)
	setString(params, "qintitle", p.QInTitle)
	setString(params, "searchIn", p.SearchIn)
	setString(params, "sources", p.Sources)
	setString(params, "domains", p.Domains)
	setString(params, "excludeDomains", p.ExcludeDomains)
	setString(params, "from", p.From)
	setString(params, "to", p.To)
	setString(params, "language", p.Language)
	setString(params, "sortBy", p.SortBy)
	setInt(params, "pageSize", p.PageSize)
	setInt(params, "page", p.Page)

	body, err := c.get("/everything", params)
	if err != nil {
		return nil, err
	}

	var result ArticlesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if result.Status == "error" {
		return nil, apiError(result.Code, result.Message)
	}
	return &result, nil
}

func (c *Client) GetSources(p SourcesParams) (*SourcesResponse, error) {
	params := url.Values{}
	setString(params, "category", p.Category)
	setString(params, "language", p.Language)
	setString(params, "country", p.Country)

	body, err := c.get("/top-headlines/sources", params)
	if err != nil {
		return nil, err
	}

	var result SourcesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if result.Status == "error" {
		return nil, apiError(result.Code, result.Message)
	}
	return &result, nil
}

// apiError maps newsapi.org error codes to descriptive messages.
func apiError(code, message string) error {
	hints := map[string]string{
		"apiKeyDisabled":    "your API key has been disabled",
		"apiKeyExhausted":   "your API key has no requests remaining",
		"apiKeyInvalid":     "your API key is invalid — check NEWSAPI_KEY",
		"apiKeyMissing":     "API key is missing",
		"parameterInvalid":  "one or more parameters are invalid",
		"parametersMissing": "required parameters are missing",
		"rateLimited":       "rate limit exceeded — slow down requests",
		"sourcesTooMany":    "too many sources requested (max 20)",
		"sourceDoesNotExist": "one or more sources do not exist",
		"unexpectedError":   "unexpected server error",
	}
	if hint, ok := hints[code]; ok {
		return fmt.Errorf("%s: %s", hint, message)
	}
	return fmt.Errorf("[%s] %s", code, message)
}

func setString(v url.Values, key, val string) {
	if val != "" {
		v.Set(key, val)
	}
}

func setInt(v url.Values, key string, val int) {
	if val > 0 {
		v.Set(key, strconv.Itoa(val))
	}
}
