package templates

var Credentials = `
AWS_ACCESS_KEY_ID = "{{.AwsAccessKeyID}}"
AWS_SECRET_ACCESS_KEY = "{{.AwsSecretKey}}"
AWS_SSH_KEY_NAME = "{{.AwsAccessSSHKey}}"
AWS_DEFAULT_REGION = "{{.AwsDefaultRegion}}"
`
