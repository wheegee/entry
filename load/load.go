package load

import (
	"fmt"
	"os"
	"os/exec"
)

func Exec(envSlice []string, args []string) error {
	command := exec.Command(args[0], args[1:]...)
	command.Env = envSlice
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func Stdout(envSlice []string) {
	for _, kv := range envSlice {
		fmt.Printf("export %s\n", kv)
	}
}
