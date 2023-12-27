package kv

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMClient interface {
	GetParametersByPath(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

type Parameters struct {
	Client SSMClient
}

func NewClient(client SSMClient) *Parameters {
	return &Parameters{
		Client: client,
	}
}

func (p *Parameters) Get(prefixes []string) (*KV, error) {
	var results *ssm.GetParametersByPathOutput
	for _, prefix := range prefixes {
		if err := validatePrefixPath(prefix); err != nil {
			return nil, err
		}

		var err error
		results, err = p.Client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path:           &prefix,
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return nil, fmt.Errorf("error retrieving parameters under path %s: %w", prefix, err)
		}
	}

	kv := make(map[string]any)
	for _, result := range results.Parameters {
		if json.Valid([]byte(*result.Value)) {
			var v map[string]any
			if err := json.Unmarshal([]byte(*result.Value), &v); err != nil {
				return nil, fmt.Errorf("error unmarshalling parameter %s: %w", *result.Name, err)
			}

			for k, v := range v {
				kv[k] = v
			}
		} else {
			k := strings.Split(*result.Name, "/")
			kv[k[len(k)-1]] = *result.Value
		}
	}

	return &KV{
		Map:  kv,
		Envs: toEnvs(kv),
	}, nil
}

func (p *Parameters) Unmarshal(prefix string, v any) error {
	withDecryption := true
	input := &ssm.GetParameterInput{
		Name:           &prefix,
		WithDecryption: &withDecryption,
	}

	result, err := p.Client.GetParameter(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("error retrieving parameter %s: %w", prefix, err)
	}

	if err := json.Unmarshal([]byte(*result.Parameter.Value), &v); err != nil {
		return fmt.Errorf("error unmarshalling parameter %s: %w", prefix, err)
	}

	return nil
}

func validatePrefixPath(prefix string) error {
	if strings.HasSuffix(prefix, "/") {
		return fmt.Errorf("prefixes must not end with a slash: %s", prefix)
	}

	if strings.Contains(prefix, "//") {
		return fmt.Errorf("prefixes must not contain double slashes: %s", prefix)
	}

	return nil
}
