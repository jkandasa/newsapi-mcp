package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Valid values sourced from https://newsapi.org/docs/endpoints
var (
	langCodes = []string{"ar", "de", "en", "es", "fr", "he", "it", "nl", "no", "pt", "ru", "sv", "ud", "zh"}

	countryCodes = []string{
		"ae", "ar", "at", "au", "be", "bg", "br", "ca", "ch", "cn",
		"co", "cu", "cz", "de", "eg", "fr", "gb", "gr", "hk", "hu",
		"id", "ie", "il", "in", "it", "jp", "kr", "lt", "lv", "ma",
		"mx", "my", "ng", "nl", "no", "nz", "ph", "pl", "pt", "ro",
		"rs", "ru", "sa", "se", "sg", "si", "sk", "th", "tr", "tw",
		"ua", "us", "ve", "za",
	}

	categories = []string{"business", "entertainment", "general", "health", "science", "sports", "technology"}
)

func registerTools(s *server.MCPServer, c *Client) {
	addTopHeadlines(s, c)
	addSearchEverything(s, c)
	addGetSources(s, c)
}

func addTopHeadlines(s *server.MCPServer, c *Client) {
	tool := mcp.NewTool("get_top_headlines",
		mcp.WithDescription("Get breaking news headlines. At least one of q, qintitle, sources, or country must be provided. 'sources' cannot be combined with 'country' or 'category'."),
		mcp.WithString("q",
			mcp.Description("Keywords or phrases to search for in the headline and body of articles."),
		),
		mcp.WithString("qintitle",
			mcp.Description("Keywords or phrases to search for in headline titles only."),
		),
		mcp.WithString("sources",
			mcp.Description("Comma-separated news source IDs (e.g. 'bbc-news,cnn'). Cannot be combined with 'country' or 'category'. Use get_sources to discover valid IDs."),
		),
		mcp.WithString("country",
			mcp.Description("2-letter ISO 3166-1 country code. Cannot be combined with 'sources'."),
			mcp.Enum(countryCodes...),
		),
		mcp.WithString("category",
			mcp.Description("News category. Cannot be combined with 'sources'."),
			mcp.Enum(categories...),
		),
		mcp.WithString("language",
			mcp.Description("2-letter ISO-639-1 language code to filter articles."),
			mcp.Enum(langCodes...),
		),
		mcp.WithNumber("page_size",
			mcp.Description("Number of results per page, 1–100 (default 20)."),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default 1)."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p := TopHeadlinesParams{
			Country:  req.GetString("country", ""),
			Category: req.GetString("category", ""),
			Sources:  req.GetString("sources", ""),
			Q:        req.GetString("q", ""),
			QInTitle: req.GetString("qintitle", ""),
			Language: req.GetString("language", ""),
			PageSize: req.GetInt("page_size", 0),
			Page:     req.GetInt("page", 0),
		}

		result, err := c.GetTopHeadlines(p)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return toJSON(result)
	})
}

func addSearchEverything(s *server.MCPServer, c *Client) {
	tool := mcp.NewTool("search_everything",
		mcp.WithDescription("Search through millions of articles from 150,000+ sources going back 5 years. At least one of q, qintitle, sources, or domains must be provided."),
		mcp.WithString("q",
			mcp.Description("Keywords or phrases (max 500 chars). Supports AND, OR, NOT operators, exact phrases in quotes, and field prefixes title:, description:, content:."),
		),
		mcp.WithString("qintitle",
			mcp.Description("Keywords or phrases to search in article titles only (max 500 chars)."),
		),
		mcp.WithString("search_in",
			mcp.Description("Comma-separated fields to restrict the q search to: title, description, content."),
		),
		mcp.WithString("sources",
			mcp.Description("Comma-separated source IDs to restrict results (max 20). Use get_sources to discover valid IDs."),
		),
		mcp.WithString("domains",
			mcp.Description("Comma-separated domains to restrict results (e.g. 'bbc.co.uk,techcrunch.com')."),
		),
		mcp.WithString("exclude_domains",
			mcp.Description("Comma-separated domains to exclude from results."),
		),
		mcp.WithString("from",
			mcp.Description("Oldest article date in ISO 8601 format (e.g. '2024-06-01' or '2024-06-01T12:00:00'). Free tier: last month only."),
		),
		mcp.WithString("to",
			mcp.Description("Newest article date in ISO 8601 format."),
		),
		mcp.WithString("language",
			mcp.Description("2-letter ISO-639-1 language code to filter articles."),
			mcp.Enum(langCodes...),
		),
		mcp.WithString("sort_by",
			mcp.Description("Sort order: relevancy (best match first), popularity (popular sources first), publishedAt (newest first, default)."),
			mcp.Enum("relevancy", "popularity", "publishedAt"),
		),
		mcp.WithNumber("page_size",
			mcp.Description("Number of results per page, 1–100 (default 100)."),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default 1)."),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p := EverythingParams{
			Q:              req.GetString("q", ""),
			QInTitle:       req.GetString("qintitle", ""),
			SearchIn:       req.GetString("search_in", ""),
			Sources:        req.GetString("sources", ""),
			Domains:        req.GetString("domains", ""),
			ExcludeDomains: req.GetString("exclude_domains", ""),
			From:           req.GetString("from", ""),
			To:             req.GetString("to", ""),
			Language:       req.GetString("language", ""),
			SortBy:         req.GetString("sort_by", ""),
			PageSize:       req.GetInt("page_size", 0),
			Page:           req.GetInt("page", 0),
		}

		if p.Q == "" && p.QInTitle == "" && p.Sources == "" && p.Domains == "" {
			return mcp.NewToolResultError("at least one of q, qintitle, sources, or domains is required"), nil
		}

		result, err := c.SearchEverything(p)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return toJSON(result)
	})
}

func addGetSources(s *server.MCPServer, c *Client) {
	tool := mcp.NewTool("get_sources",
		mcp.WithDescription("List available news sources on newsapi.org. Returns source IDs that can be used as filters in get_top_headlines and search_everything. All parameters are optional."),
		mcp.WithString("category",
			mcp.Description("Filter sources by news category."),
			mcp.Enum(categories...),
		),
		mcp.WithString("language",
			mcp.Description("Filter sources by 2-letter ISO-639-1 language code."),
			mcp.Enum(langCodes...),
		),
		mcp.WithString("country",
			mcp.Description("Filter sources by 2-letter ISO 3166-1 country code."),
			mcp.Enum(countryCodes...),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p := SourcesParams{
			Category: req.GetString("category", ""),
			Language: req.GetString("language", ""),
			Country:  req.GetString("country", ""),
		}

		result, err := c.GetSources(p)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return toJSON(result)
	})
}

func toJSON(v any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}
	return mcp.NewToolResultText(string(data)), nil
}
