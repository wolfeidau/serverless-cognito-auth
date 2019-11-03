package triggers

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/serverless-cognito-auth/pkg/hooks"

	"github.com/stretchr/testify/require"
)

var (
	data = []byte(`{
		"version": "123",
		"triggerSource": "PreAuthentication_Authentication",
		"region": "AWSRegion",
		"userPoolId": "string",
		"userName": "string",
		"callerContext": 
			{
				"awsSdkVersion": "string",
				"clientId": "string"
			},
		"request":
			{
				"userAttributes": {
					"string": "string"
				}
			},
		"response": {}
	}`)
)

func Test_Handler(t *testing.T) {

	assert := require.New(t)

	ct := &CognitoTriggers{}

	hk := &hooks.LambdaHooks{}

	hk.EventReceived = func(ctx context.Context) context.Context {
		log.Info().Msg("EventReceived")
		return ctx
	}
	hk.EventSent = func(ctx context.Context) {
		log.Info().Msg("EventSent")

	}

	di := NewDispatcher(ct, hk)

	data, err := di.Handler(context.TODO(), data)

	assert.Nil(err)
	assert.NotNil(data)
}
