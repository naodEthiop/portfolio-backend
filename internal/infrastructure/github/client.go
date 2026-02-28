package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient(token string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		token:      token,
	}
}

type SearchOptions struct {
	Sort    string
	Order   string
	PerPage int
}

type APIError struct {
	StatusCode         int
	Message            string
	DocumentationURL   string
	RateLimitRemaining int
	RateLimitReset     time.Time
}

func (e *APIError) Error() string {
	if e == nil {
		return "github api error"
	}
	if e.Message != "" {
		return fmt.Sprintf("github api: status=%d message=%s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("github api: status=%d", e.StatusCode)
}

func (e *APIError) RateLimited() bool {
	if e == nil {
		return false
	}
	// GitHub commonly returns 403 for rate limiting; also guard on headers.
	if e.RateLimitRemaining == 0 && !e.RateLimitReset.IsZero() {
		return true
	}
	if e.StatusCode == http.StatusTooManyRequests {
		return true
	}
	return false
}

type repoItem struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	HTMLURL         string    `json:"html_url"`
	Description     string    `json:"description"`
	Language        string    `json:"language"`
	Homepage        string    `json:"homepage"`
	Topics          []string  `json:"topics"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	Archived        bool      `json:"archived"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
}

type SearchResponse struct {
	TotalCount int        `json:"total_count"`
	Items      []repoItem `json:"items"`
}

type errorBody struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

// SearchRepositories calls GitHub Search API for repositories matching q.
// q example: "user:naodEthiop topic:portfolio"
func (c *Client) SearchRepositories(ctx context.Context, q string, opt SearchOptions) (*SearchResponse, error) {
	u := url.URL{Scheme: "https", Host: "api.github.com", Path: "/search/repositories"}
	qv := u.Query()
	qv.Set("q", q)
	if opt.Sort != "" {
		qv.Set("sort", opt.Sort)
	}
	if opt.Order != "" {
		qv.Set("order", opt.Order)
	}
	if opt.PerPage > 0 {
		qv.Set("per_page", strconv.Itoa(opt.PerPage))
	}
	u.RawQuery = qv.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiErr := &APIError{
			StatusCode:         resp.StatusCode,
			RateLimitRemaining: parseHeaderInt(resp.Header.Get("X-RateLimit-Remaining")),
			RateLimitReset:     parseHeaderUnixTime(resp.Header.Get("X-RateLimit-Reset")),
		}

		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		var eb errorBody
		if err := json.Unmarshal(body, &eb); err == nil {
			apiErr.Message = eb.Message
			apiErr.DocumentationURL = eb.DocumentationURL
		} else if len(body) > 0 {
			apiErr.Message = string(body)
		}
		return nil, apiErr
	}

	var out SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("github api: empty response")
		}
		return nil, err
	}
	return &out, nil
}

func parseHeaderInt(value string) int {
	if value == "" {
		return -1
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return -1
	}
	return n
}

func parseHeaderUnixTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	sec, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(sec, 0).UTC()
}
