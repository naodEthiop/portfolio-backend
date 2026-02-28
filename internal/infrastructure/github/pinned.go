package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type pinnedRepo struct {
	ID              int64     `json:"databaseId"`
	Name            string    `json:"name"`
	NameWithOwner   string    `json:"nameWithOwner"`
	URL             string    `json:"url"`
	Description     string    `json:"description"`
	HomepageURL     string    `json:"homepageUrl"`
	StargazerCount  int       `json:"stargazerCount"`
	ForkCount       int       `json:"forkCount"`
	IsArchived      bool      `json:"isArchived"`
	UpdatedAt       time.Time `json:"updatedAt"`
	PushedAt        time.Time `json:"pushedAt"`
	PrimaryLanguage *struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
	RepositoryTopics struct {
		Nodes []struct {
			Topic struct {
				Name string `json:"name"`
			} `json:"topic"`
		} `json:"nodes"`
	} `json:"repositoryTopics"`
}

type pinnedResponse struct {
	Data struct {
		User struct {
			PinnedItems struct {
				Nodes []pinnedRepo `json:"nodes"`
			} `json:"pinnedItems"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *Client) ListPinnedRepositories(ctx context.Context, login string, first int) ([]pinnedRepo, error) {
	// GraphQL endpoint is effectively unusable without a token for most deployments,
	// but we treat it as optional. If token is missing, return empty.
	if c.token == "" {
		return []pinnedRepo{}, nil
	}
	if first <= 0 || first > 50 {
		first = 6
	}

	const query = `
query($login: String!, $first: Int!) {
  user(login: $login) {
    pinnedItems(first: $first, types: [REPOSITORY]) {
      nodes {
        ... on Repository {
          databaseId
          name
          nameWithOwner
          url
          description
          homepageUrl
          stargazerCount
          forkCount
          isArchived
          updatedAt
          pushedAt
          primaryLanguage { name }
          repositoryTopics(first: 20) { nodes { topic { name } } }
        }
      }
    }
  }
}`

	body, err := json.Marshal(map[string]any{
		"query": query,
		"variables": map[string]any{
			"login": login,
			"first": first,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.github.com/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github graphql status: %s", resp.Status)
	}

	var out pinnedResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Errors) > 0 {
		return nil, fmt.Errorf("github graphql error: %s", out.Errors[0].Message)
	}

	repos := out.Data.User.PinnedItems.Nodes
	if repos == nil {
		return []pinnedRepo{}, nil
	}
	return repos, nil
}
