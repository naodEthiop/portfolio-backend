package usecase

import (
	"context"
	"sort"
	"strings"
	"time"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"
)

type PublicProjectsService struct {
	github *GithubService
	repo   repository.ProjectRepository
}

func NewPublicProjectsService(github *GithubService, repo repository.ProjectRepository) *PublicProjectsService {
	return &PublicProjectsService{github: github, repo: repo}
}

func (s *PublicProjectsService) List(ctx context.Context, user string, topic string) (*entities.PublicProjectsResponse, error) {
	manual, err := s.repo.ListPublic(ctx)
	if err != nil {
		return nil, err
	}

	githubProjects, err := s.github.ListProjects(ctx, user, topic)
	if err != nil {
		return nil, err
	}

	var pinnedGithub []entities.GithubProject
	// pinned repos should show even if they don't have the "portfolio" topic
	if pinned, pinErr := s.github.ListPinnedProjects(ctx, user, ""); pinErr == nil {
		pinnedGithub = pinned
	}
	if len(pinnedGithub) == 0 && len(githubProjects) > 0 {
		candidates := make([]entities.GithubProject, len(githubProjects))
		copy(candidates, githubProjects)
		sort.SliceStable(candidates, func(i, j int) bool {
			if candidates[i].Stars != candidates[j].Stars {
				return candidates[i].Stars > candidates[j].Stars
			}
			return candidates[i].PushedAt.After(candidates[j].PushedAt)
		})
		limit := 3
		if len(candidates) < limit {
			limit = len(candidates)
		}
		pinnedGithub = candidates[:limit]
	}

	githubByRepoURL := make(map[string]entities.GithubProject, len(githubProjects)+len(pinnedGithub))
	pinnedRepoURL := make(map[string]struct{}, len(pinnedGithub))

	for _, p := range githubProjects {
		u := normalizeURL(p.HTMLURL)
		if u == "" {
			continue
		}
		githubByRepoURL[u] = p
	}
	for _, p := range pinnedGithub {
		u := normalizeURL(p.HTMLURL)
		if u == "" {
			continue
		}
		pinnedRepoURL[u] = struct{}{}
		// ensure pinned repos are present even if they don't match the search topic
		githubByRepoURL[u] = p
	}

	consumedGithub := make(map[string]struct{}, len(githubByRepoURL))
	unified := make([]entities.PublicProject, 0, len(githubByRepoURL)+len(manual))

	for _, p := range manual {
		repoURL := normalizeURL(p.RepoURL)
		if repoURL != "" && repoURL != "#" {
			if gp, ok := githubByRepoURL[repoURL]; ok {
				_, isPinned := pinnedRepoURL[repoURL]
				pp := mapManualGithubToPublic(p, gp, isPinned || p.Featured)
				unified = append(unified, pp)
				consumedGithub[repoURL] = struct{}{}
				continue
			}
		}
		unified = append(unified, mapManualToPublic(p, p.Featured))
	}

	for repoURL, gp := range githubByRepoURL {
		if _, ok := consumedGithub[repoURL]; ok {
			continue
		}
		_, isPinned := pinnedRepoURL[repoURL]
		unified = append(unified, mapGithubToPublic(gp, isPinned))
	}

	pinned := make([]entities.PublicProject, 0, 8)
	for _, p := range unified {
		if p.Pinned {
			pinned = append(pinned, p)
		}
	}

	sort.SliceStable(pinned, func(i, j int) bool {
		return pinned[i].UpdatedAt.After(pinned[j].UpdatedAt)
	})

	sort.SliceStable(unified, func(i, j int) bool {
		if unified[i].Pinned != unified[j].Pinned {
			return unified[i].Pinned
		}
		return unified[i].UpdatedAt.After(unified[j].UpdatedAt)
	})

	return &entities.PublicProjectsResponse{
		Pinned:   pinned,
		Projects: unified,
	}, nil
}

func mapGithubToPublic(p entities.GithubProject, pinned bool) entities.PublicProject {
	updatedAt := p.PushedAt
	if updatedAt.IsZero() {
		updatedAt = p.UpdatedAt
	}
	if updatedAt.IsZero() {
		updatedAt = time.Now().UTC()
	}

	tech := make([]string, 0, len(p.Topics)+1)
	for _, t := range p.Topics {
		if strings.TrimSpace(t) == "" {
			continue
		}
		tech = append(tech, t)
	}

	return entities.PublicProject{
		ID:          "github:" + int64ToString(p.ID),
		Source:      entities.ProjectSourceGithub,
		Name:        p.Name,
		Description: strings.TrimSpace(p.Description),
		TechStack:   tech,
		Topics:      sliceOrEmpty(p.Topics),
		Language:    strings.TrimSpace(p.Language),
		Stars:       p.Stars,
		Forks:       p.Forks,
		DemoURL:     strings.TrimSpace(p.Homepage),
		RepoURL:     strings.TrimSpace(p.HTMLURL),
		Featured:    false,
		Pinned:      pinned,
		UpdatedAt:   updatedAt,
	}
}

func mapManualToPublic(p entities.Project, pinned bool) entities.PublicProject {
	return entities.PublicProject{
		ID:          p.ID.String(),
		Source:      entities.ProjectSourceManual,
		Name:        p.Title,
		Description: strings.TrimSpace(firstNonEmpty(p.ShortDescription, p.Description)),
		TechStack:   sliceOrEmpty([]string(p.TechStack)),
		Topics:      []string{},
		Language:    "",
		Stars:       0,
		Forks:       0,
		DemoURL:     strings.TrimSpace(p.DemoURL),
		RepoURL:     strings.TrimSpace(p.RepoURL),
		Featured:    p.Featured,
		Pinned:      pinned,
		UpdatedAt:   p.UpdatedAt,
	}
}

func mapManualGithubToPublic(manual entities.Project, gh entities.GithubProject, pinned bool) entities.PublicProject {
	updatedAt := gh.PushedAt
	if updatedAt.IsZero() {
		updatedAt = gh.UpdatedAt
	}
	if updatedAt.IsZero() {
		updatedAt = manual.UpdatedAt
	}

	return entities.PublicProject{
		ID:          manual.ID.String(),
		Source:      entities.ProjectSourceManual,
		Name:        manual.Title,
		Description: strings.TrimSpace(firstNonEmpty(manual.ShortDescription, manual.Description, gh.Description)),
		TechStack:   sliceOrEmpty([]string(manual.TechStack)),
		Topics:      sliceOrEmpty(gh.Topics),
		Language:    strings.TrimSpace(gh.Language),
		Stars:       gh.Stars,
		Forks:       gh.Forks,
		DemoURL:     strings.TrimSpace(firstNonEmpty(manual.DemoURL, gh.Homepage)),
		RepoURL:     strings.TrimSpace(firstNonEmpty(gh.HTMLURL, manual.RepoURL)),
		Featured:    manual.Featured,
		Pinned:      pinned,
		UpdatedAt:   updatedAt,
	}
}

func normalizeURL(value string) string {
	return strings.TrimSuffix(strings.TrimSpace(value), "/")
}

func int64ToString(v int64) string {
	// avoid importing strconv in a tight mapping file; keep local.
	// strconv is fine; but this keeps dependencies minimal.
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [32]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + (v % 10))
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func sliceOrEmpty(in []string) []string {
	if in == nil {
		return []string{}
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
