# entry

Entry is a binary for use with the `AWS_LAMBDA_EXEC_WRAPPER`. It has a very narrow scope of functionality that does not negotiate.

## Requisites

### 1. Data Layer

You store your environment variables in encrypted SSM keys as JSON.

`ssm://path/to/envjson` 

```json
{
  "ENVAR_0": "value_0",
  "ENVAR_1": "value_1"
}
```

### 2. Runtime Layer

Your execution context has AWS credentials available to the default credential chain that are sufficient for accessing the data layer (SSM read && KMS decrypt).

```json
{
    "Sid": "AllowSSMParameterAccess",
    "Effect": "Allow",
    "Action": [
        "ssm:GetParametersByPath",
        "ssm:GetParameter",
        "ssm:GetParameters",
        "ssm:PutParameter",
        "kms:Decrypt",
        "kms:Encrypt"
    ],
    "Resource": [
        "arn:aws:ssm:{{AWS_REGION}}:{{AWS_ACCOUNT_ID}}:parameter/path/to/envjson",
        "arn:aws:ssm:{{AWS_REGION}}:{{AWS_ACCOUNT_ID}}:parameter/path/to/envjson/*"
    ]
}
```

## Usage

### Core Pattern
```
# eval export statements
eval $(entry --path /path/to/envjson)
```

```
# inect environment variables into child process env
go run cmd/main.go --path /path/to/env/json -- env
```

```
# merge mutliple ssm env paths
go run cmd/main.go --path /path/to/envjson_1 --path /path/to/envjson_2
```

### AWS Lambda
1. Build your lambda image with `entry` at `/opt/entry`
2. Deploy your lambda with the envar: `AWS_LAMBDA_EXEC_WRAPPER=/opt/entry --path /path/to/envjson -- ${@}`
