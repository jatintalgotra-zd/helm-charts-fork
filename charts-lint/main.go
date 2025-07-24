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

	chartDirectories := make([]string, 0)

	for _, chart := range chartsDir {
		if chart.IsDir() {
			dir := filepath.Join("charts", chart.Name())
			err = helmDependencyUpdate(dir)
			if err != nil {
				fmt.Println(err)
			}
			chartDirectories = append(chartDirectories, dir)
		}
	}

	lintingResult := helmLint(chartDirectories)

	for _, err = range lintingResult.Messages {
		fmt.Println(err)
	}
}
