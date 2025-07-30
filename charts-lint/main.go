package main

import (
	"errors"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"os"
	"path/filepath"

	githubpkg "charts-lint/github"
	helmpkg "charts-lint/helm"
)

var ErrMissingBaseRef = errors.New("GITHUB_BASE_REF environment variable not set")

// helmChart holds chart name and its path.
type helmChart struct {
	name string
	path string
}

func main() {
	// Get charts changed in the current PR
	changedCharts, err := githubpkg.GetDiff()
	if err != nil {
		fmt.Println(err)
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

		os.Exit(1)
	}
}

// getDiff returns a list of unique chart names under the 'charts/' directory that have been modified between the base branch and HEAD.
//func getDiff() ([]string, error) {
//	base := os.Getenv("GITHUB_BASE_REF")
//
//	if base == "" {
//		return nil, ErrMissingBaseRef
//	}
//
//	arg := fmt.Sprintf("origin/%s...HEAD", base)
//
//	// run git diff to get differences in base branch and head
//	cmd := exec.Command("git", "diff", "--name-only", arg)
//	cmd.Stderr = os.Stderr
//
//	out, err := cmd.Output()
//	if err != nil {
//		return nil, err
//	}
//
//	changed := strings.Split(string(out), "\n")
//	check := make(map[string]bool)
//	result := make([]string, 0)
//
//	for _, line := range changed {
//		split := strings.Split(line, "/")
//		if split[0] == "charts" && len(split) > 1 {
//			// Target only paths like charts/<chart>
//			dir := split[1]
//			if !check[dir] {
//				check[dir] = true
//
//				result = append(result, dir)
//			}
//		}
//	}
//
//	return result, nil
//}
