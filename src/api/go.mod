module github.com/jacksonwendel/nebulaos/src/api

go 1.25.0

require (
	github.com/gophercloud/gophercloud v1.14.1
	github.com/jacksonwendel/nebulaos/src/providers v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.11.2
	github.com/nats-io/nats.go v1.49.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.41.2 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.293.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.7 // indirect
	github.com/aws/smithy-go v1.24.1 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/jacksonwendel/nebulaos/src/providers => ../providers
