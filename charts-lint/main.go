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

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "Chart.yaml" {
			dir := filepath.Dir(path)

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

		return nil
	})
	if err != nil {
		fmt.Println(err)
	} else if len(lintingErrors) > 0 {
		for _, e := range lintingErrors {
			fmt.Printf("Linting error in chart %v: %v \n", e.Chart, e.Err)
		}
	}
}
