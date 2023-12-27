# entry

[![Main](https://github.com/linecard/entry/actions/workflows/main.yaml/badge.svg)](https://github.com/linecard/entry/actions/workflows/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/linecard/entry)](https://goreportcard.com/report/github.com/linecard/entry)

A package and execution wrapper for AWS SSM Parameter expansion.

## Usage

### CLI

```bash
Usage: main [-g] [-p PREFIX] [-v] [COMMAND [ARGUMENTS [ARGUMENTS ...]]]

Positional arguments:
  COMMAND                Command to run
  ARGUMENTS              Command arguments

Options:
  -g                     Do not inherit environment
  -p PREFIX              SSM prefixes to source
  -v                     Verbose output
  --help, -h             display this help and exit
  --version              display version and exit
```

### Package

```json
{
  "foo": "oof",
  "bar": "rab"
}
```

```go
var data struct {
  Foo string `json:"foo"`
  Bar string `json:"bar"`
}

awsConfig, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
  panic("error loading AWS credentials")
}

p := &kv.Parameters{Client: ssm.NewFromConfig(awsConfig)}

// Get parameters under provided prefixes
params, err := p.Get([]string{"/dev"})
if err != nil {
  panic("error getting parameters")
}

// Unmarshal specific parameter into provided data type
if err := p.Unmarshal("/dev/foobar", &data); err != nil {
  panic("error unmarshalling parameter")
}
```
