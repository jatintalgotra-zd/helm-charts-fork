package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type lintError struct {
	Chart string
	Err   error
}

func execCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	lintingErrors := make([]lintError, 0)

	chartsDir, err := os.ReadDir("charts")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, chart := range chartsDir {
		if !chart.IsDir() {
			continue
		}

		dir := filepath.Join("charts", chart.Name())
		if _, err = os.Stat(dir + "/Chart.yaml"); err != nil {
			fmt.Println("Chart not found in directory: ", dir, "\nError: ", err)
			return
		}

		fmt.Println("Updating dependencies for: ", dir)
		err = execCommand("helm", "dependency", "update", dir)
		if err != nil {
			lintingErrors = append(lintingErrors, lintError{
				Chart: dir,
				Err:   err,
			})
		}

		fmt.Println("Linting chart: ", dir)
		err = execCommand("helm", "lint", dir)
		if err != nil {
			lintingErrors = append(lintingErrors, lintError{
				Chart: dir,
				Err:   err,
			})
		}
	}

	if err != nil {
		fmt.Println(err)
	} else if len(lintingErrors) > 0 {
		for _, e := range lintingErrors {
			fmt.Printf("Linting error in chart %v: %v \n", e.Chart, e.Err)
		}
		os.Exit(1)
	}
}
