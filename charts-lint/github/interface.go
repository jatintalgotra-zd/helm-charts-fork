package github

import (
	"context"
	gh "github.com/google/go-github/v74/github"
)

type Client interface {
	ListFiles(ctx context.Context, owner, repo string, number int, opts *gh.ListOptions) ([]*gh.CommitFile, *gh.Response, error)
}
