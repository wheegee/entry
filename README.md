# entry

Entry provides a simple solution for managing application configuration during any phase of the development lifecycle.

## Usage

### CLI
```shell
./entry -p /path/to/env -- env
```

### Dockerfile
```Dockerfile
FROM scratch
COPY --from=ghcr.io/entry/entry:latest /ko-app/entry /opt/entry

ENTRYPOINT ["/opt/entry", "-p", "/path/to/env", "--"] 
CMD ["env"]
```

## Storing Environment

Entry assumes the usage of AWS SSM as the backing data store for your environments.

### SSM Parameter
1. The parameter type shall be of type Secret String.
2. The parameter value shall be of JSON format.
3. The parameter JSON schema shall be of the form...

```json
{
    "ENVAR_1": "value_1",
    "ENVAR_2": "value_2"
}
```

### Caller Permissions

1. The caller shall have AWS credentials available to the [credential provider chain](https://docs.aws.amazon.com/sdkref/latest/guide/standardized-credentials.html#credentialProviderChain).
2. The caller shall have permissions akin to the following...

```json
{
    "sid": "ssmAccess",
    "effect": "Allow",
    "action": [
        "ssm:GetParameter",
        "kms:Decrypt"
    ],
    "resource": [
        "arn:aws:ssm:${AWS_ACCOUNT_REGION}:${AWS_ACCOUNT_ID}:parameter/*"
    ]
}
```