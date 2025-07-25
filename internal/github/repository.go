package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v64/github"
)

// Repository represents a GitHub repository with its owner and name.
// It contains a GitHub client and context for making API calls.
type Repository struct {
	ctx    context.Context
	client *github.Client
	Owner  string
	Repo   string
}

// NewRepository creates a new instance of Repository.
// It requires the repository owner, repository name, and a GitHub client.
func NewRepository(owner, repo string, client *github.Client) *Repository {
	return &Repository{
		Owner:  owner,
		Repo:   repo,
		client: client,
		ctx:    context.Background(),
	}
}

// LatestRelease retrieves the latest release for the repository.
func (r *Repository) LatestRelease() (*Release, error) {
	repositoryRelease, _, err := r.client.Repositories.GetLatestRelease(r.ctx, r.Owner, r.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}

	release := &Release{}
	if err := release.FromRepositoryRelease(repositoryRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}

// GetRelease retrieves a specific release for the repository based on the provided tag.
func (r *Repository) GetRelease(tag string) (*Release, error) {
	repositoryRelease, _, err := r.client.Repositories.GetReleaseByTag(r.ctx, r.Owner, r.Repo, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets for release tag %q: %w", tag, err)
	}

	release := &Release{}
	if err := release.FromRepositoryRelease(repositoryRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}

// LatestIncludingPreRelease retrieves the most recently published release for the repository,
// including pre-releases. This returns the newest release by published date, regardless of
// whether it's a regular release or pre-release.
func (r *Repository) LatestIncludingPreRelease(perPage int) (*Release, error) {
	// List all releases including pre-releases
	opts := &github.ListOptions{
		PerPage: perPage,
	}

	repositoryReleases, _, err := r.client.Repositories.ListReleases(r.ctx, r.Owner, r.Repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	if len(repositoryReleases) == 0 {
		return nil, fmt.Errorf("no releases found for %s/%s", r.Owner, r.Repo)
	}

	// Find the most recent release by published date
	var latestRelease *github.RepositoryRelease
	for i, release := range repositoryReleases {
		if i == 0 || release.PublishedAt == nil || latestRelease.PublishedAt == nil {
			latestRelease = release

			continue
		}

		// Compare the timestamps - need to use the Time property of Timestamp
		if release.PublishedAt.After(latestRelease.PublishedAt.Time) {
			latestRelease = release
		}
	}

	// Convert to our Release type
	release := &Release{}
	if err := release.FromRepositoryRelease(latestRelease); err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	return release, nil
}
