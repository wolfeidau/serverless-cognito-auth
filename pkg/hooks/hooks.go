package hooks

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
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

func NewDefaultServerHooks() *LambdaHooks {

	hooks := &LambdaHooks{}

	hooks.EventReceived = func(ctx context.Context) context.Context {
		log.Info().Msg("EventReceived")

		return context.WithValue(ctx, startTimeKey, time.Now())
	}

	hooks.EventSent = func(ctx context.Context) {

		defer libhoney.Flush()

		startTime := getStartTime(ctx)

		log.Info().Time("startTime", startTime).Msg("EventSent")

		lc, _ := lambdacontext.FromContext(ctx)

		ev := libhoney.NewEvent()
		ev.Add(map[string]interface{}{
			"function_name":    lambdacontext.FunctionName,
			"function_version": lambdacontext.FunctionVersion,
			"request_id":       lc.AwsRequestID,
			"duration_ms":      time.Since(startTime).Milliseconds(),
		})
		ev.Send()

		log.Info().Msg("libhoney EventSent")

	}

	return hooks
}

func getStartTime(ctx context.Context) time.Time {
	return ctx.Value(startTimeKey).(time.Time)
}
