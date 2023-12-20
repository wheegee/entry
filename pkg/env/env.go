package env

import (
	"os"
	"os/exec"

	"github.com/linecard/entry/pkg/kv"
)

type Env struct {
	Parameters *kv.Parameters
	Envs       []string
	Command    string
	Arguments  []string
	Pristine   bool
}

func Configure(client kv.SSMClient, command string, arguments ...string) Env {
	return Env{
		Parameters: &kv.Parameters{
			Client: client,
		},
		Command:   command,
		Arguments: arguments,
		Pristine:  false,
	}
}

func (e *Env) Execute() error {
	cmd := exec.Command(e.Command, e.Arguments...)

	if e.Pristine {
		cmd.Env = e.Envs
	} else {
		cmd.Env = append(os.Environ(), e.Envs...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
