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

func inLambda() bool {
	_, exists := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	return exists
}

func main() {
	preDash, hasDash, postDash := extract.Argv(os.Args)

	paths, verbose, err := extract.ParseFlags(preDash)
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	mergedParams, err := extract.SSM(ctx, ssmClient, paths)
	if err != nil {
		log.Fatalf("Failed to fetch SSM parameters: %v", err)
	}

	envSlice, err := transform.ToEnvSlice(mergedParams)
	if err != nil {
		log.Fatalf("Failed to transform SSM parameters: %v", err)
	}

	if hasDash && len(postDash) > 0 {
		if verbose {
			load.LogMasked(envSlice)
		}

		mergedWithParent := append(os.Environ(), envSlice...)
		if err := load.Exec(mergedWithParent, postDash); err != nil {
			log.Fatalf("Failed child execution: %s\n%v", strings.Join(postDash, " "), err)
		}
		return
	}

	// don't print export statements to stdout if in lamdba.
	if !inLambda() {
		load.Stdout(envSlice)
		return
	}

	log.Fatalf("When running in a Lambda, a command must be provided using the '--' separator.")
}
