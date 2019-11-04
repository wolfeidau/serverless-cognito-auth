package hooks

import (
	"context"
	"os"

	libhoney "github.com/honeycombio/libhoney-go"
	"github.com/rs/zerolog/log"
)

type userCtxKeyType string

const startTimeKey userCtxKeyType = "startTime"

// LambdaHooks used to instrument lambda event handlers
type LambdaHooks struct {
	EventReceived func(context.Context) context.Context
	EventSent     func(context.Context)
}

func Init() error {

	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	dataset := os.Getenv("HONEYCOMB_DATASET")

	if apiKey == "" || dataset == "" {
		log.Info().Msg("no honeycomb configuration found")
		return nil
	}

	log.Info().Str("Dataset", dataset).Msg("honeycomb init")

	return libhoney.Init(libhoney.Config{
		WriteKey: apiKey,
		Dataset:  dataset,
	})
}
