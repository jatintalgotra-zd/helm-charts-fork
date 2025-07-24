package main

import (
	"fmt"
	"os"
	"path/filepath"

	helm "helm.sh/helm/v3/pkg/action"
)

//type lintError struct {
//	Chart string
//	Err   error
//}

func helmDependencyUpdate(path string) error {

}

func helmLint(paths []string) *helm.LintResult {
	lint := helm.NewLint()
	result := lint.Run(paths, nil)
	return result
}

func main() {
	//lintingErrors := make([]lintError, 0)

	chartsDir, err := os.ReadDir("charts")
	if err != nil {
		fmt.Println(err)
		return
	}

	chartDirectories := make([]string, 0)

	for _, chart := range chartsDir {
		if chart.IsDir() {
			chartDirectories = append(chartDirectories, filepath.Join("charts", chart.Name()))
		}
	}

	lintingResult := helmLint(chartDirectories)

	for _, err = range lintingResult.Messages {
		fmt.Println(err)
	}
}
