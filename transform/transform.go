package transform

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// ToEnvSlice transforms a slice of SSM parameters into a slice of environment variable strings.
func ToEnvSlice(mergedParams []types.Parameter) (envSlice []string, err error) {
	for _, param := range mergedParams {
		var jsonEnv map[string]string
		if err := json.Unmarshal([]byte(*param.Value), &jsonEnv); err != nil {
			return nil, err
		}

		for key, value := range jsonEnv {
			envSlice = append(envSlice, fmt.Sprintf("%s=%v", key, value))
		}
	}

	return envSlice, nil
}