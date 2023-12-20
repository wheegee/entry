package env

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/linecard/entry/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	var (
		prefixes      = []string{"/user-api/postgres", "/admin-api/redis"}
		command       = "printenv"
		arguments     = []string{}
		ssmParameters = ssm.GetParametersOutput{
			Parameters: []types.Parameter{
				{
					Name:  aws.String("/user-api/postgres"),
					Value: aws.String(`{"postgres_user":"user-api","postgres_pass":"password"}`),
				},
				{
					Name:  aws.String("/admin-api/redis"),
					Value: aws.String(`{"redis_host":"redis.internal","redis_port":"6379"}`),
				},
			},
		}
	)
	cases := []struct {
		name          string
		prefixes      []string
		command       string
		arguments     []string
		pristine      bool
		ssmParameters ssm.GetParametersOutput
		awsError      error
		expectedError error
	}{
		{
			name:          "success executing with env",
			prefixes:      prefixes,
			command:       command,
			arguments:     arguments,
			pristine:      false,
			ssmParameters: ssmParameters,
			awsError:      nil,
			expectedError: nil,
		},
		{
			name:          "success executing with pristine env",
			prefixes:      prefixes,
			command:       command,
			arguments:     arguments,
			pristine:      true,
			ssmParameters: ssmParameters,
			awsError:      nil,
			expectedError: nil,
		},
		{
			name:          "error executing with env",
			prefixes:      prefixes,
			command:       command,
			arguments:     arguments,
			pristine:      false,
			ssmParameters: ssm.GetParametersOutput{},
			awsError:      nil,
			expectedError: fmt.Errorf("ahh!"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := mock.MockSSMClient{}
			client.GetParametersFunc = func(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
				return &c.ssmParameters, c.awsError
			}

			e := Configure(&client, c.command, c.arguments...)
			params, err := e.Parameters.Get(c.prefixes)
			if err != nil {
				assert.Equal(t, c.expectedError, err)
			}
			e.Envs = params.Envs
			e.Pristine = c.pristine

			output, err := captureStdout(e.Execute)
			if err != nil {
				assert.Equal(t, c.expectedError, err)
			} else if c.pristine {
				// Concat params.Envs into a single, multinline string, and assert against output
				var envs string
				for _, env := range params.Envs {
					envs += env + "\n"
				}

				assert.Equal(t, envs, output)
			} else {
				for _, env := range params.Envs {
					assert.Contains(t, output, env)
				}
			}
		})
	}
}

func captureStdout(f func() error) (string, error) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	w.Close()
	os.Stdout = origStdout
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}
