package main

import (
	"context"
	"os"

	"github.com/linecard/entry/extract"
	"github.com/linecard/entry/internal/util"
	"github.com/linecard/entry/load"
	"github.com/linecard/entry/transform"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/rs/zerolog/log"
)

var ctx context.Context
var ssmClient *ssm.Client

func init() {
	ctx = context.Background()

	util.SetLogLevel()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to load AWS SDK config")
	}

	ssmClient = ssm.NewFromConfig(awsConfig)
}

func main() {
	preDash, hasDash, postDash := extract.Argv(os.Args)

	paths, err := extract.ParseFlags(preDash)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse flags")
	}

	mergedParams, err := extract.SSM(ctx, ssmClient, paths)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to fetch SSM parameters")
	}

	envSlice, err := transform.ToEnvSlice(mergedParams)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to transform SSM parameters")
	}

	if hasDash && len(postDash) > 0 {
		mergedWithParent := append(os.Environ(), envSlice...)
		if err := load.Exec(mergedWithParent, postDash); err != nil {
			log.Fatal().Err(err).Strs("cmd", postDash).Msgf("Failed child process execution")
		}
		return
	}
}
