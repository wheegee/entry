package kv

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMClient interface {
	GetParameters(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
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
	withDecryption := true
	input := &ssm.GetParametersInput{
		Names:          prefixes,
		WithDecryption: &withDecryption,
	}

	results, err := p.Client.GetParameters(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("error retrieving parameters: %w", err)
	}

	kv := make(map[string]any)
	for _, result := range results.Parameters {
		if json.Valid([]byte(*result.Value)) {
			var v map[string]any
			if err := json.Unmarshal([]byte(*result.Value), &v); err != nil {
				return nil, fmt.Errorf("error unmarshalling parameter: %w", err)
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
		return fmt.Errorf("error retrieving parameter: %w", err)
	}

	if err := json.Unmarshal([]byte(*result.Parameter.Value), &v); err != nil {
		return fmt.Errorf("error unmarshalling parameter: %w", err)
	}

	return nil
}
