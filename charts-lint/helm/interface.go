package helm

import (
	"io"

	"helm.sh/helm/v3/pkg/action"
)

type dependency interface {
	List(string, io.Writer) error
}

type manager interface {
	Update() error
}

type linter interface {
	Run(paths []string, vals map[string]any) *action.LintResult
}
