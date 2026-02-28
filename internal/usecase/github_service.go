package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/infrastructure/github"

	"golang.org/x/sync/singleflight"
)

type GithubService struct {
	client  *github.Client
	ttl     time.Duration
	perPage int

	mu    sync.RWMutex
	cache map[string]githubCacheEntry
	sf    singleflight.Group
}

type githubCacheEntry struct {
	expires  time.Time
	projects []entities.GithubProject
}

func NewGithubService(client *github.Client, ttl time.Duration, perPage int) *GithubService {
	if perPage <= 0 || perPage > 100 {
		perPage = 30
	}
	return &GithubService{
		client:  client,
		ttl:     ttl,
		perPage: perPage,
		cache:   make(map[string]githubCacheEntry),
	}
}

func (s *GithubService) ListProjects(ctx context.Context, user string, topic string) ([]entities.GithubProject, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("%w: github client not configured", ErrUpstream)
	}

	key := user + "|" + topic
	if s.ttl > 0 {
		if cached, ok := s.getCached(key); ok {
			return cached, nil
		}
	}

	val, err, _ := s.sf.Do(key, func() (any, error) {
		if s.ttl > 0 {
			if cached, ok := s.getCached(key); ok {
				return cached, nil
			}
		}

		q := fmt.Sprintf("user:%s topic:%s", user, topic)
		resp, err := s.client.SearchRepositories(ctx, q, github.SearchOptions{
			Sort:    "updated",
			Order:   "desc",
			PerPage: s.perPage,
		})
		if err != nil {
			var apiErr *github.APIError
			if errors.As(err, &apiErr) && apiErr.RateLimited() {
				return nil, fmt.Errorf("%w: %s", ErrRateLimited, apiErr.Error())
			}
			return nil, fmt.Errorf("%w: %s", ErrUpstream, err.Error())
		}

		projects := make([]entities.GithubProject, 0, len(resp.Items))
		for _, it := range resp.Items {
			topics := it.Topics
			if topics == nil {
				topics = []string{}
			}
			projects = append(projects, entities.GithubProject{
				ID:          it.ID,
				Name:        it.Name,
				FullName:    it.FullName,
				HTMLURL:     it.HTMLURL,
				Homepage:    it.Homepage,
				Description: it.Description,
				Topics:      topics,
				Language:    it.Language,
				Stars:       it.StargazersCount,
				Forks:       it.ForksCount,
				Archived:    it.Archived,
				UpdatedAt:   it.UpdatedAt,
				PushedAt:    it.PushedAt,
			})
		}

		sort.Slice(projects, func(i, j int) bool {
			return projects[i].PushedAt.After(projects[j].PushedAt)
		})

		if s.ttl > 0 {
			s.mu.Lock()
			s.cache[key] = githubCacheEntry{
				expires:  time.Now().Add(s.ttl),
				projects: projects,
			}
			s.mu.Unlock()
		}

		out := make([]entities.GithubProject, len(projects))
		copy(out, projects)
		return out, nil
	})
	if err != nil {
		return nil, err
	}
	projects, ok := val.([]entities.GithubProject)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected github service result", ErrUpstream)
	}
	return projects, nil
}

func (s *GithubService) getCached(key string) ([]entities.GithubProject, bool) {
	s.mu.RLock()
	entry, ok := s.cache[key]
	s.mu.RUnlock()
	if !ok || time.Now().After(entry.expires) {
		return nil, false
	}
	out := make([]entities.GithubProject, len(entry.projects))
	copy(out, entry.projects)
	return out, true
}

func (s *GithubService) ListPinnedProjects(ctx context.Context, user string, topic string) ([]entities.GithubProject, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("%w: github client not configured", ErrUpstream)
	}

	key := "pinned|" + user + "|" + topic
	if s.ttl > 0 {
		if cached, ok := s.getCached(key); ok {
			return cached, nil
		}
	}

	val, err, _ := s.sf.Do(key, func() (any, error) {
		if s.ttl > 0 {
			if cached, ok := s.getCached(key); ok {
				return cached, nil
			}
		}

		repos, err := s.client.ListPinnedRepositories(ctx, user, 12)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrUpstream, err.Error())
		}

		projects := make([]entities.GithubProject, 0, len(repos))
		for _, r := range repos {
			topics := make([]string, 0, len(r.RepositoryTopics.Nodes))
			for _, n := range r.RepositoryTopics.Nodes {
				if n.Topic.Name != "" {
					topics = append(topics, n.Topic.Name)
				}
			}
			if topic != "" && !containsString(topics, topic) {
				continue
			}
			lang := ""
			if r.PrimaryLanguage != nil {
				lang = r.PrimaryLanguage.Name
			}
			projects = append(projects, entities.GithubProject{
				ID:          r.ID,
				Name:        r.Name,
				FullName:    r.NameWithOwner,
				HTMLURL:     r.URL,
				Homepage:    r.HomepageURL,
				Description: r.Description,
				Topics:      topics,
				Language:    lang,
				Stars:       r.StargazerCount,
				Forks:       r.ForkCount,
				Archived:    r.IsArchived,
				UpdatedAt:   r.UpdatedAt,
				PushedAt:    r.PushedAt,
			})
		}

		sort.Slice(projects, func(i, j int) bool {
			return projects[i].PushedAt.After(projects[j].PushedAt)
		})

		if s.ttl > 0 {
			s.mu.Lock()
			s.cache[key] = githubCacheEntry{
				expires:  time.Now().Add(s.ttl),
				projects: projects,
			}
			s.mu.Unlock()
		}

		out := make([]entities.GithubProject, len(projects))
		copy(out, projects)
		return out, nil
	})
	if err != nil {
		return nil, err
	}
	projects, ok := val.([]entities.GithubProject)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected github service result", ErrUpstream)
	}
	return projects, nil
}

func containsString(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}
