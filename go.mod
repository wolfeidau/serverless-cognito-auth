module github.com/wolfeidau/serverless-cognito-auth

go 1.13

require (
	github.com/DataDog/zstd v1.4.1 // indirect
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.25.25
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/limitgroup v0.0.0-20150612190941-6abd8d71ec01 // indirect
	github.com/facebookgo/muster v0.0.0-20150708232844-fd3d7953fd52 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/honeycombio/libhoney-go v1.12.1
	github.com/klauspost/compress v1.9.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rs/zerolog v1.15.0
	github.com/stretchr/testify v1.4.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/alexcesaro/statsd.v2 v2.0.0 // indirect
)

replace github.com/aws/aws-lambda-go => github.com/wolfeidau/aws-lambda-go v1.13.3-0.20191023124854-b5b7267d297d
