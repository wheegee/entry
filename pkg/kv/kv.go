package kv

import "strconv"

type KV struct {
	Map  map[string]any
	Envs []string
}

func toEnvs(kv map[string]any) []string {
	env := make([]string, 0, len(kv))
	for k, v := range kv {
		switch value := v.(type) {
		case bool:
			v = strconv.FormatBool(value)
		case float64:
			v = strconv.FormatFloat(value, 'f', -1, 64)
		}
		env = append(env, k+"="+v.(string))
	}

	return env
}
