package main

import (
	"fmt"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func helmDependencyUpdate(path string) error {
	settings := cli.New()

	manager := &downloader.Manager{
		Out:              os.Stdout,
		ChartPath:        path,
		Getters:          getter.All(settings),
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
	}

	if err := manager.Update(); err != nil {
		return fmt.Errorf("failed to update dependencies for %s: %w", path, err)
	}
	return nil
}

func helmLint(paths []string) *action.LintResult {
	lint := action.NewLint()
	result := lint.Run(paths, nil)
	return result
}

func main() {
	chartsDir, err := os.ReadDir("charts")
	if err != nil {
		fmt.Println(err)
		return
	}

	failedCharts := make([]string, 0)

	for _, chart := range chartsDir {
		if chart.IsDir() {
			dir := filepath.Join("charts", chart.Name())
			fmt.Printf("\n=== Processing Chart: %s ===\n", chart.Name())

			fmt.Printf("-> Updating dependencies...\n")
			err = helmDependencyUpdate(dir)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("OK: Dependencies updated successfully or already up-to-date.")
			}

			fmt.Printf("\n-> Running lint...\n")
			result := helmLint([]string{dir})

			if len(result.Errors) > 0 {
				failedCharts = append(failedCharts, chart.Name())
				fmt.Println(result.Errors[0])
			} else {
				fmt.Printf("OK: Lint succeeded.\n")
			}
		}
	}

	if len(failedCharts) > 0 {
		fmt.Printf("\n=== Failed Charts: %v\n", failedCharts)
		os.Exit(1)
	}
}
