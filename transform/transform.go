package transform

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// ToEnvSlice transforms a slice of SSM parameters into a slice of environment variable strings.
func ToEnvSlice(mergedParams []types.Parameter, verbose bool) (envSlice []string, err error) {
	for _, param := range mergedParams {
		var jsonEnv map[string]string
		if err := json.Unmarshal([]byte(*param.Value), &jsonEnv); err != nil {
			return nil, err
		}

		for key, value := range jsonEnv {
			if verbose {
				maskedValue := strings.Repeat("*", len(value))
				log.Printf("export %s=%s\n", key, maskedValue)
			}
			envSlice = append(envSlice, fmt.Sprintf("%s=%v", key, value))
		}
	}

	return envSlice, nil
}
