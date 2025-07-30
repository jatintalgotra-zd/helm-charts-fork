package utils

import (
	"charts-lint/models"
	"errors"
	"testing"
)

func TestGetGithubClient(t *testing.T) {
	testcases := []struct {
		name        string
		setEnv      func()
		expectedErr error
	}{
		{
			name: "Github Token Not Set",
			setEnv: func() {
			},
			expectedErr: ErrMissingToken,
		},
		{
			name: "Github Token Set",
			setEnv: func() {
				t.Setenv("GITHUB_TOKEN", "test-token")
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setEnv()

			_, err := GetGithubClient()
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("GetGithubClient() error = %v, wantErr %v", err, tc.expectedErr)
			}
		})
	}
}

func TestSummary(t *testing.T) {
	testcases := []struct {
		name   string
		failed []models.HelmChart
		passed []models.HelmChart
	}{
		{
			name: "TestSummary",
			failed: []models.HelmChart{
				{},
			},
			passed: []models.HelmChart{
				{},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			Summary(tc.failed, tc.passed)
		})
	}
}
