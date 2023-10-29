package main

import (
	"github.com/alexflint/go-arg"
	"github.com/linecard/entry/pkg/env"
)

type args struct {
	Pristine      bool     `arg:"-g,--" help:"Do not inherit environment"`
	ParamPrefixes []string `arg:"-p,--,separate" placeholder:"PREFIX" help:"SSM prefixes to source"`
	Command       string   `arg:"positional" help:"Command to run"`
	Arguments     []string `arg:"positional" help:"Command arguments"`
}

func main() {
	var args args
	arg.MustParse(&args)

	if len(args.ParamPrefixes) == 0 {
		panic("no prefixes specified")
	}

	e := env.Configure(args.Command, args.Arguments...)
	params, err := e.Parameters.Get(args.ParamPrefixes)
	if err != nil {
		panic(err)
	}
	e.Envs = params.Envs

	if err := e.Execute(); err != nil {
		panic(err)
	}
}
