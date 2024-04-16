package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/linecard/entry/extract"
	"github.com/linecard/entry/load"
	"github.com/linecard/entry/transform"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var ctx context.Context
var ssmClient *ssm.Client

func init() {
	ctx = context.Background()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config, %v", err)
	}

	ssmClient = ssm.NewFromConfig(awsConfig)
}

func main() {
	preDash, postDash := extract.Argv(os.Args)

	paths, verbose, err := extract.ParseFlags(preDash)
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	mergedParams, err := extract.SSM(ctx, ssmClient, paths)
	if err != nil {
		log.Fatalf("Failed to fetch SSM parameters: %v", err)
	}

	envSlice, err := transform.ToEnvSlice(mergedParams, verbose)
	if err != nil {
		log.Fatalf("Failed to transform SSM parameters: %v", err)
	}

	if len(postDash) > 0 {
		mergedEnv := append(os.Environ(), envSlice...)
		if err := load.Exec(mergedEnv, postDash); err != nil {
			log.Fatalf("Failed execution: %s\n%v", strings.Join(postDash, " "), err)
		}
		return
	}

	load.Stdout(envSlice)
}
