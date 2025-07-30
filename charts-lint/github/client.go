package github

import (
	"context"
	"os"
	"strconv"
	"strings"
)

type githubClient struct {
	pr PullRequest
}

func New(pr PullRequest) *githubClient {
	return &githubClient{pr: pr}
}

func (c *githubClient) GetDiff() ([]string, error) {
	owner := os.Getenv("REPOSITORY_OWNER")
	repoPath := os.Getenv("REPOSITORY_NAME")
	repo := strings.Split(repoPath, "/")[1]

	prNumber, err := strconv.Atoi(os.Getenv("PR_NUMBER"))
	if err != nil {
		return nil, err
	}

	commitFiles, _, err := c.pr.ListFiles(context.Background(), owner, repo, prNumber, nil)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	check := make(map[string]bool)

	for _, file := range commitFiles {
		name := file.GetFilename()

		split := strings.Split(name, "/")
		if split[0] == "charts" && len(split) > 1 {
			// Target only paths like charts/<chart>
			dir := split[1]
			if !check[dir] {
				check[dir] = true

				files = append(files, dir)
			}
		}
	}

	return files, nil
}
