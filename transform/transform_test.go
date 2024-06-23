package transform

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToEnvSlice(t *testing.T) {
	tests := []struct {
		name        string
		inputParams []types.Parameter
		expectedEnv []string
		expectErr   bool
	}{
		{
			name: "single valid parameter",
			inputParams: []types.Parameter{
				{
					Name:  aws.String("/path/to/env"),
					Value: aws.String(`{"DB_HOST": "localhost", "DB_PORT": "5432"}`),
				},
			},
			expectedEnv: []string{"DB_HOST=localhost", "DB_PORT=5432"},
			expectErr:   false,
		},
		{
			name: "invalid JSON format",
			inputParams: []types.Parameter{
				{
					Name:  aws.String("/path/to/env"),
					Value: aws.String(`{"DB_HOST": "localhost", "DB_PORT": "5432"`), // Missing closing brace
				},
			},
			expectedEnv: nil,
			expectErr:   true,
		},
		{
			name:        "empty parameters",
			inputParams: []types.Parameter{},
			expectedEnv: []string{},
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envSlice, err := ToEnvSlice(tt.inputParams)
			if tt.expectErr {
				require.Error(t, err, "Test %s should have returned an error", tt.name)
			} else {
				require.NoError(t, err, "Test %s should not have returned an error", tt.name)
				assert.ElementsMatch(t, tt.expectedEnv, envSlice, "Test %s failed: expected envSlice %v, got %v", tt.name, tt.expectedEnv, envSlice)
			}
		})
	}
}
