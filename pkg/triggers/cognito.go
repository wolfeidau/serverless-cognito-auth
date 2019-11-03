package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	// ErrInvalidDomain reject email domain
	ErrInvalidDomain = errors.New("invalid email domain")
	// ErrInvalidSignUpCode reject supplied invite code
	ErrInvalidSignUpCode = errors.New("invalid invite code, please check the value")
	// ErrMissingEmail email missing from event attributes
	ErrMissingEmail = errors.New("missing email")
)

// CognitoPoolEventEnvelope cognito events for a pool with dynamic request/response to cater for
// different types
type CognitoPoolEventEnvelope struct {
	Version       string                                    `json:"version"`
	TriggerSource string                                    `json:"triggerSource"`
	Region        string                                    `json:"region"`
	UserPoolID    string                                    `json:"userPoolId"`
	CallerContext events.CognitoEventUserPoolsCallerContext `json:"callerContext"`
	UserName      string                                    `json:"userName"`
	Request       *json.RawMessage                          `json:"request"`
	Response      *json.RawMessage                          `json:"response"`
}

// CognitoTriggers handles all the cognito triggers supported in this service
type CognitoTriggers struct {
}

// PreSignUp pre sign up used to implement checks prior before signup
func (ct *CognitoTriggers) PreSignUp(ctx context.Context, evt *CognitoPoolEventEnvelope, request *events.CognitoEventUserPoolsPreSignupRequest) (*events.CognitoEventUserPoolsPreSignupResponse, error) {
	log.Info().Str("UserPoolID", evt.UserPoolID).Str("UserName", evt.UserName).Msg("PreSignUp")

	if signUpInviteCode := os.Getenv("SIGNUP_INVITE_CODE"); signUpInviteCode != "" {
		inviteCode, ok := request.UserAttributes["invite_code"]
		if !ok {
			log.Error().Msg("user attributes missing invite_code")
			return nil, ErrInvalidSignUpCode
		}

		if signUpInviteCode != inviteCode {
			log.Error().Str("InviteCodeSupplied", inviteCode).Str("ExpectedInviteCode", signUpInviteCode).Msg("invite_code didn't match configured value")
			return nil, ErrInvalidSignUpCode
		}

	}

	// if WHITELIST_DOMAIN is configured then verify the email is from that domain
	if whitelistDomain := os.Getenv("WHITELIST_DOMAIN"); whitelistDomain != "" {
		email, ok := request.UserAttributes["email"]
		if !ok {
			log.Error().Msg("user attributes missing email")
			return nil, ErrInvalidDomain
		}
		tokens := strings.Split(email, "@")
		if len(tokens) != 2 {
			log.Error().Str("email", email).Msg("missing email domain")
			return nil, ErrInvalidDomain
		}
		if tokens[1] != whitelistDomain {
			log.Error().Str("email", email).Str("whitelistDomain", whitelistDomain).Msg("email domain doesn't match whitelist domain")
			return nil, ErrInvalidDomain
		}

	}

	return &events.CognitoEventUserPoolsPreSignupResponse{}, nil
}

// PreAuthentication pre authentication
func (ct *CognitoTriggers) PreAuthentication(ctx context.Context, evt *CognitoPoolEventEnvelope, request *events.CognitoEventUserPoolsPreAuthenticationRequest) (*events.CognitoEventUserPoolsPreAuthenticationResponse, error) {
	log.Info().Str("UserPoolID", evt.UserPoolID).Str("UserName", evt.UserName).Msg("PreAuthentication")
	return &events.CognitoEventUserPoolsPreAuthenticationResponse{}, nil
}

// PreToken pretoken lambda used to override roles/groups/claims
func (ct *CognitoTriggers) PreToken(ctx context.Context, evt *CognitoPoolEventEnvelope, request *events.CognitoEventUserPoolsPreTokenGenRequest) (*events.CognitoEventUserPoolsPreTokenGenResponse, error) {
	log.Info().Str("UserPoolID", evt.UserPoolID).Str("UserName", evt.UserName).Msg("PreToken")
	return &events.CognitoEventUserPoolsPreTokenGenResponse{}, nil
}

// PostConfirmationSignUp post confirmation used to notify of signups
func (ct *CognitoTriggers) PostConfirmationSignUp(ctx context.Context, evt *CognitoPoolEventEnvelope, request *events.CognitoEventUserPoolsPostConfirmationRequest) (*events.CognitoEventUserPoolsPostConfirmationResponse, error) {
	log.Info().Str("UserPoolID", evt.UserPoolID).Str("UserName", evt.UserName).Msg("PostConfirmationSignUp")

	email, ok := request.UserAttributes["email"]
	if !ok {
		log.Error().Msg("user attributes missing email")
		return nil, ErrMissingEmail
	}

	sess := session.Must(session.NewSession())

	svc := sns.New(sess)

	res, err := svc.Publish(&sns.PublishInput{
		TopicArn: aws.String(os.Getenv("SIGNUP_SNS_TOPIC")),
		Subject:  aws.String("new signup"),
		Message:  aws.String(fmt.Sprintf("signup of user %s", email)),
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to publish PostConfirmationSignUp message")
		return nil, err
	}

	log.Info().Str("MessageId", aws.StringValue(res.MessageId)).Msg("PostConfirmationSignUp message published")

	return &events.CognitoEventUserPoolsPostConfirmationResponse{}, nil
}
