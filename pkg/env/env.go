package env

import (
	"context"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/linecard/entry/pkg/kv"
)

type Env struct {
	Parameters *kv.Parameters
	Envs       []string
	Command    string
	Arguments  []string
	Pristine   bool
}

func Configure(command string, arguments ...string) Env {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("error loading AWS credentials")
	}

	return Env{
		Parameters: &kv.Parameters{
			Client: ssm.NewFromConfig(awsConfig),
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
