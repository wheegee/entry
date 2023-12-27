package kv

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/linecard/entry/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var (
		prefixes = []string{"/user-api"}
		kv       = KV{
			Map: map[string]any{
				"postgres_user": "user-api",
				"postgres_pass": "password",
				"redis_host":    "redis.internal",
				"redis_port":    "6379",
			},
			Envs: []string{
				"postgres_user=user-api",
				"postgres_pass=password",
				"redis_host=redis.internal",
				"redis_port=6379",
			},
		}
		ssmParameters = ssm.GetParametersByPathOutput{
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
		kv            KV
		ssmParameters ssm.GetParametersByPathOutput
		awsError      error
		expectedError error
	}{
		{
			name:          "success retrieving parameters",
			prefixes:      prefixes,
			kv:            kv,
			ssmParameters: ssmParameters,
			awsError:      nil,
			expectedError: nil,
		},
		{
			name:          "error retrieving parameters",
			prefixes:      prefixes,
			kv:            KV{},
			ssmParameters: ssm.GetParametersByPathOutput{},
			awsError:      fmt.Errorf("ahh!"),
			expectedError: fmt.Errorf("error retrieving parameters under path %s: %w", prefixes[0], fmt.Errorf("ahh!")),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := mock.MockSSMClient{}
			client.GetParametersByPathFunc = func(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error) {
				return &c.ssmParameters, c.awsError
			}

			p := NewClient(&client)

			params, err := p.Get(c.prefixes)
			if err != nil {
				assert.Equal(t, c.expectedError, err)
			} else {
				assert.Equal(t, c.kv.Map, params.Map)
				assert.ElementsMatch(t, c.kv.Envs, params.Envs)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type (
		postgresSecrets struct {
			PostgresUser string `json:"postgres_user"`
			PostgresPass string `json:"postgres_pass"`
		}
	)
	var (
		prefix       = "/user-api/postgres"
		ssmParameter = ssm.GetParameterOutput{
			Parameter: &types.Parameter{
				Name:  aws.String("/user-api/postgres"),
				Value: aws.String(`{"postgres_user":"user-api","postgres_pass":"password"}`),
			},
		}
	)
	cases := []struct {
		name          string
		prefix        string
		targetStruct  postgresSecrets
		ssmParameter  ssm.GetParameterOutput
		awsError      error
		expectedError error
	}{
		{
			name:          "success unmarshalling parameter",
			prefix:        prefix,
			targetStruct:  postgresSecrets{},
			ssmParameter:  ssmParameter,
			awsError:      nil,
			expectedError: nil,
		},
		{
			name:          "error retrieving parameter",
			prefix:        prefix,
			targetStruct:  postgresSecrets{},
			ssmParameter:  ssm.GetParameterOutput{},
			awsError:      fmt.Errorf("ahh!"),
			expectedError: fmt.Errorf("error retrieving parameter %s: %w", prefix, fmt.Errorf("ahh!")),
		},
		{
			name:          "error unmarshalling parameter",
			prefix:        prefix,
			targetStruct:  postgresSecrets{},
			ssmParameter:  ssmParameter,
			awsError:      nil,
			expectedError: fmt.Errorf("ahh!"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := mock.MockSSMClient{}
			client.GetParameterFunc = func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
				return &c.ssmParameter, c.awsError
			}

			p := NewClient(&client)

			err := p.Unmarshal(c.prefix, &c.targetStruct)
			if err != nil {
				assert.Equal(t, c.expectedError, err)
			} else {
				assert.Equal(t, "user-api", c.targetStruct.PostgresUser)
				assert.Equal(t, "password", c.targetStruct.PostgresPass)
			}
		})
	}
}
