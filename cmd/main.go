package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"project-root/env"
	"project-root/ssm"
)

func main() {
	cfg, ssmClient, err := ssm.SetupClient()
	if err != nil {
		log.Fatalf("Unable to setup SSM client, %v", err)
	}

	ssmFlag := flag.NewFlagSet("example", flag.ExitOnError)
	var ssms []string
	ssmFlag.Func("ssm", "Specify SSM parameters to fetch", func(s string) error {
		ssms = append(ssms, s)
		return nil
	})

	ssmFlag.Parse(os.Args[1:])
	args := ssmFlag.Args()

	var additionalArgs []string
	for i, arg := range args {
		if arg == "--" {
			additionalArgs = args[i+1:]
			args = args[:i]
			break
		}
	}

	if len(ssms) > 0 {
		params, err := ssm.FetchParameters(ssmClient, ssms)
		if err != nil {
			log.Fatalf("Failed to fetch SSM parameters: %v", err)
		}
		for _, value := range params {
			if err := env.SetEnvVarsFromJSON(value); err != nil {
				log.Fatalf("Failed to set env vars: %v", err)
			}
		}
	}

	if len(additionalArgs) > 0 {
		cmd := exec.Command(additionalArgs[0], additionalArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
		}
	} else {
		fmt.Println("No command specified to execute.")
	}
}
