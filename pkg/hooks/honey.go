package hooks

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	libhoney "github.com/honeycombio/libhoney-go"
	"github.com/rs/zerolog/log"
)

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
