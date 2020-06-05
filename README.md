# Gales Dips Propane Report Download

### Description
Creates an xlsx report file in S3, then returns a signed URL for file

## Dependancies
Note: this would change after converting to go mod
``` bash
$ dep ensure -add github.com/aws/aws-sdk-go/service
$ dep ensure -add github.com/machinebox/graphql
```