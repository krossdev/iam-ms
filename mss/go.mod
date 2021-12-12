module github.com/krossdev/iam-ms/mss

go 1.17

replace github.com/krossdev/iam-ms/msc => ../msc

require (
	github.com/airbrake/gobrake/v5 v5.1.1
	github.com/blang/semver/v4 v4.0.0
	github.com/fatih/color v1.13.0
	github.com/krossdev/iam-ms/msc v0.0.0-00010101000000-000000000000
	github.com/mitchellh/mapstructure v1.4.3
	github.com/nats-io/nats.go v1.13.1-0.20211122170419-d7c1d78a50fc
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/caio/go-tdigest v3.1.0+incompatible // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/minio/highwayhash v1.0.1 // indirect
	github.com/nats-io/jwt/v2 v2.2.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
)
