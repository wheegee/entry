package transform

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/rs/zerolog/log"
)

// ToEnvSlice transforms a slice of SSM parameters into a slice of environment variable strings.
func ToEnvSlice(mergedParams []types.Parameter) (envSlice []string, err error) {
	for _, param := range mergedParams {
		var path string
		var keys []string
		var jsonEnv map[string]string

		path = *param.Name

		if err := json.Unmarshal([]byte(*param.Value), &jsonEnv); err != nil {
			return nil, err
		}

		for key, value := range jsonEnv {
			keys = append(keys, key)
			envSlice = append(envSlice, fmt.Sprintf("%s=%v", key, value))
		}

		log.Info().Str("path", path).Strs("keys", keys).Msg("parameters loaded")
	}

	return envSlice, nil
}
