package main

import (
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/serverless-cognito-auth/pkg/hooks"
	"github.com/wolfeidau/serverless-cognito-auth/pkg/triggers"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Time("start", time.Now()).Msg("triggers")

	err := hooks.Init()
	if err != nil {
		log.Fatal().Msg("failed to init hooks")
	}

	mon := hooks.NewDefaultServerHooks()
	ct := triggers.NewCognitoTriggers()
	dispatcher := triggers.NewDispatcher(ct, mon)

	lambda.Start(dispatcher.Handler)
}
