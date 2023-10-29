package kv

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Parameters struct {
	Client *ssm.Client
}

func (p *Parameters) Get(prefixes []string) (*KV, error) {
	withDecryption := true
	input := &ssm.GetParametersInput{
		Names:          prefixes,
		WithDecryption: &withDecryption,
	}

	results, err := p.Client.GetParameters(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	kv := make(map[string]any)
	for _, result := range results.Parameters {
		if json.Valid([]byte(*result.Value)) {
			var v map[string]any
			if err := json.Unmarshal([]byte(*result.Value), &v); err != nil {
				return nil, err
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
		return err
	}

	if err := json.Unmarshal([]byte(*result.Parameter.Value), &v); err != nil {
		return err
	}

	return nil
}
