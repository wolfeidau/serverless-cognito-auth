package eventstream

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/rs/zerolog/log"
)

type EventPayload struct {
	AwsRequestID string            `json:"aws_request_id,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`
}

type Publisher struct {
	svc snsiface.SNSAPI
}

// NewPublisher new publisher
func NewPublisher(cfgs ...*aws.Config) *Publisher {
	sess := session.Must(session.NewSession(cfgs...))
	svc := sns.New(sess)
	return &Publisher{svc: svc}
}

func (pb *Publisher) SendEvent(ctx context.Context, topicName string, subject string, requestId string, attributes map[string]string) error {

	payload := &EventPayload{
		AwsRequestID: requestId,
		Attributes:   attributes,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	res, err := pb.svc.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicName),
		Subject:  aws.String(subject),
		Message:  aws.String(string(data)),
	})
	if err != nil {
		return err
	}

	log.Info().Str("id", aws.StringValue(res.MessageId)).Msg("event sent")

	return nil
}

func (pb *Publisher) SendNotification(ctx context.Context, topicName string, subject string, msg string) error {

	res, err := pb.svc.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicName),
		Subject:  aws.String(subject),
		Message:  aws.String(msg),
	})
	if err != nil {
		return err
	}

	log.Info().Str("id", aws.StringValue(res.MessageId)).Msg("message sent")

	return nil
}
