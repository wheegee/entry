package main

import (
	"context"

	"github.com/alexflint/go-arg"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/linecard/entry/pkg/env"
)

var (
	version = "dev"
	commit  = "none"
)

type args struct {
	Pristine      bool     `arg:"-g,--" help:"Do not inherit environment"`
	ParamPrefixes []string `arg:"-p,--,separate" placeholder:"PREFIX" help:"SSM prefixes to source"`
	Command       string   `arg:"positional" help:"Command to run"`
	Arguments     []string `arg:"positional" help:"Command arguments"`
}

func (args) Version() string {
	return "entry " + version + " " + commit
}

func main() {
	var args args
	arg.MustParse(&args)
	if len(args.ParamPrefixes) == 0 {
		panic("no prefixes specified")
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("error loading AWS credentials")
	}

	e := env.Configure(ssm.NewFromConfig(awsConfig), args.Command, args.Arguments...)
	params, err := e.Parameters.Get(args.ParamPrefixes)
	if err != nil {
		panic(err)
	}
	e.Envs = params.Envs

	if err := e.Execute(); err != nil {
		panic(err)
	}
}
