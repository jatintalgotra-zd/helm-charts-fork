package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v74/github"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"

	githubpkg "charts-lint/github"
	helmpkg "charts-lint/helm"
)

// ErrMissingToken is returned when the GITHUB_TOKEN environment variable is not set.
var ErrMissingToken = errors.New("GITHUB_TOKEN environment variable not set")

// helmChart holds chart name and its path.
type helmChart struct {
	name string
	path string
}

// getGithubClient helper function to get GitHub client from access token.
func getGithubClient() (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, ErrMissingToken
	}

	return github.NewClient(nil).WithAuthToken(token), nil
}

// summary helper function to print summary of helm lint.
func summary(failedCharts, passedCharts []helmChart) {
	// Passed charts
	if len(passedCharts) > 0 {
		fmt.Printf("\nTOTAL CHARTS UPDATED: %v\n", len(passedCharts)+len(failedCharts))
		fmt.Println("\n===> Charts that passed helm lint:")

		for _, chart := range passedCharts {
			fmt.Printf("-> %s - path: %s\n", chart.name, chart.path)
		}
	}

	// Failed charts
	if len(failedCharts) > 0 {
		fmt.Println("\n===> Charts not passing helm lint:")

		for _, chart := range failedCharts {
			fmt.Printf("-> %s - path: %s\n", chart.name, chart.path)
		}
	}
}

func main() {
	// get client for GitHub
	client, err := getGithubClient()
	if err != nil {
		panic(err)
	}

	// dependency injection for GitHub client
	c := githubpkg.New(client.PullRequests)

	// Get charts changed in the current PR
	changedCharts, err := c.GetDiff()
	if err != nil {
		panic(err)
	}

	// early exit for no updated charts
	if len(changedCharts) == 0 {
		fmt.Println("No charts updated.")
		return
	}

	failedCharts := make([]helmChart, 0)
	passedCharts := make([]helmChart, 0)

	// Process each changed chart
	for _, chart := range changedCharts {
		dir := filepath.Join("charts", chart)

		dep := action.NewDependency()
		settings := cli.New()
		manager := &downloader.Manager{
			Out:              os.Stdout,
			ChartPath:        dir,
			Getters:          getter.All(settings),
			RepositoryConfig: settings.RepositoryConfig,
			RepositoryCache:  settings.RepositoryCache,
		}
		lint := action.NewLint()

		// dependency injection
		helm := helmpkg.New(dep, manager, lint)

		hc := helmChart{name: chart, path: dir}
		fmt.Printf("\n=== Processing Chart: %s ===\n", chart)

		// Step 1: Update dependencies
		err = helm.DependencyUpdate(dir)
		if err != nil {
			failedCharts = append(failedCharts, hc)

			fmt.Println(err)

			continue
		}

		// Step 2: Lint
		errorMessages := helm.Lint([]string{dir})
		if errorMessages != nil {
			failedCharts = append(failedCharts, hc)

			for _, err = range errorMessages {
				fmt.Println(err)
			}
		} else {
			passedCharts = append(passedCharts, hc)

			fmt.Printf("OK: Lint succeeded.\n")
		}
	}

	// Final summary
	summary(failedCharts, passedCharts)

	if len(failedCharts) > 0 {
		os.Exit(1)
	}
}
