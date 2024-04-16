package load

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Exec executes a command with the provided environment variables.
func Exec(envSlice []string, args []string) error {
	command := exec.Command(args[0], args[1:]...)
	command.Env = envSlice
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

// Stdout prints the environment variables to stdout.
func Stdout(envSlice []string) {
	for _, kv := range envSlice {
		fmt.Printf("export %s\n", kv)
	}
}

func LogMasked(envSlice []string) {
	for _, kv := range envSlice {
		parts := strings.SplitN(kv, "=", 2)
		unmaskedKey := parts[0]
		maksedValue := strings.Repeat("*", len(parts[1]))
		log.Printf("export %s=%s\n", unmaskedKey, maksedValue)
	}
}
