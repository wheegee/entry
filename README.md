# entry

Entry is a convention for defining a containers environment variables via SSM.

## Usage

### Dockerfile
```Dockerfile
FROM ghcr.io/entry/entry:0.7.2 as entry

FROM golang:1.22.2 as build
# Build your application

FROM scratch
COPY --from=entry /ko-entry/entry /opt/entry
COPY --from=build /dist/app /var/task/app

ENTRYPOINT /opt/entry -p /path/to/json/env -- /var/task/app  
```

### CLI
```shell
# 1. Print env export statements to stdout.
./entry -p /path/to/json/env

# 2. Export env to current shell.
eval $(./entry -p /path/to/json/env)

# 3. Execute child process with the env.
./entry -p /path/to/json/env -- env

# 4. Merge multiple envs.
./entry -p /path/to/json/env1 -p /path/to/json/env2 -- env
```

## Requisites

Assuming you are storing your environment at `ssm://path/to/json/env`...

### SSM Parameter
1. The parameter type shall be of secret string.
2. The parameter value shall be of JSON format.
3. The parameter JSON value schema shall be of the form...

```json
{
    "FOO": "bar",
    "BAZ": "faz"
}
```

### Permissions

This is the gist, need to improve the example.

```json
{
    "sid": "ssmAccess",
    "effect": "Allow",
    "action": [
        "ssm:GetParameter",
        "kms:Decrypt"
    ],
    "resource": [
        "arn:aws:ssm:us-west-2:123456789012:parameter/path/to/env/json"
    ]
}
```