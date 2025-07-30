package helm

import (
	"bytes"
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/lint/support"
)

type helm struct {
	dep     dependency
	manager manager
	lint    linter
}

func New(dep dependency, manager manager, lint linter) *helm {
	return &helm{dep: dep, manager: manager, lint: lint}
}

// DependencyUpdate runs 'helm dependency update' for a given chart path.
func (h *helm) DependencyUpdate(path string) error {
	var buff bytes.Buffer
	// runs helm dependency list
	err := h.dep.List(path, &buff)
	if err != nil {
		return err
	}

	// early exit for no dependencies warning
	if strings.HasPrefix(buff.String(), "WARNING") {
		return nil
	}

	fmt.Printf("-> Updating dependencies for %v...\n", path)

	// runs helm dependency update
	if err = h.manager.Update(); err != nil {
		return fmt.Errorf("failed to update dependencies for %s: %w", path, err)
	}

	fmt.Printf("-> Successfully updated dependencies for %v\n\n", path)

	return nil
}

// Lint runs 'helm lint' on the given chart paths.
func (h *helm) Lint(paths []string) []error {
	fmt.Printf("-> Running lint for %v...\n", paths[0])

	result := h.lint.Run(paths, nil)
	errors := make([]error, 0)

	if !action.HasWarningsOrErrors(result) {
		return nil
	}

	for _, msg := range result.Messages {
		if msg.Severity == support.ErrorSev {
			errors = append(errors, msg)
		} else if msg.Severity == support.WarningSev {
			fmt.Println(msg)
		}
	}

	return errors
}
