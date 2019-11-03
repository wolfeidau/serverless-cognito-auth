package triggers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/wolfeidau/serverless-cognito-auth/pkg/hooks"
)

// Cognito trigger sources
const (
	PreSignUp                             = "PreSignUp_SignUp"
	PreSignUpAdminCreateUser              = "PreSignUp_AdminCreateUser"
	PostConfirmationConfirmSignUp         = "PostConfirmation_ConfirmSignUp"
	PostConfirmationConfirmForgotPassword = "PostConfirmation_ConfirmForgotPassword"
	PreAuthentication                     = "PreAuthentication_Authentication"
	PostAuthentication                    = "PostAuthentication_Authentication"
	PreTokenHostedAuth                    = "TokenGeneration_HostedAuth"
	PreTokenAuthentication                = "TokenGeneration_Authentication"
	PreTokenNewPasswordChallenge          = "TokenGeneration_NewPasswordChallenge"
	PreTokenAuthenticateDevice            = "TokenGeneration_AuthenticateDevice"
	PreTokenRefreshTokens                 = "TokenGeneration_RefreshTokens"
)

// Dispatcher trigger dispatcher
type Dispatcher struct {
	ct  *CognitoTriggers
	mon *hooks.LambdaHooks
}

// NewDispatcher create a new trigger dispacher
func NewDispatcher(ct *CognitoTriggers, mon *hooks.LambdaHooks) *Dispatcher {
	return &Dispatcher{ct: ct, mon: mon}
}

// Handler dispatches events to cognito triggers
func (di *Dispatcher) Handler(ctx context.Context, payload json.RawMessage) (json.RawMessage, error) {

	ctx = di.mon.EventReceived(ctx)

	defer di.mon.EventSent(ctx)

	evt := new(CognitoPoolEventEnvelope)

	err := json.Unmarshal(payload, evt)
	if err != nil {
		return nil, err
	}

	switch evt.TriggerSource {
	case PreSignUp:
		err := di.PreSignUp(ctx, evt)
		if err != nil {
			return nil, err
		}
	case PreAuthentication:
		err := di.PreAuthentication(ctx, evt)
		if err != nil {
			return nil, err
		}
	case PostConfirmationConfirmSignUp:
		err := di.PostConfirmationSignUp(ctx, evt)
		if err != nil {
			return nil, err
		}
	case PreTokenHostedAuth, PreTokenAuthentication, PreTokenNewPasswordChallenge, PreTokenAuthenticateDevice, PreTokenRefreshTokens:
		err := di.PreAuthentication(ctx, evt)
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(evt)
}

// PreSignUp presign up trigger dispatcher
func (di *Dispatcher) PreSignUp(ctx context.Context, evt *CognitoPoolEventEnvelope) error {

	req := new(events.CognitoEventUserPoolsPreSignupRequest)
	err := json.Unmarshal(*evt.Request, req)
	if err != nil {
		return err
	}

	res, err := di.ct.PreSignUp(ctx, evt, req)
	if err != nil {
		return err
	}

	rawRes, err := RawResponse(res)
	if err != nil {
		return err
	}

	evt.Response = rawRes

	return nil
}

// PreAuthentication preauthentication trigger dispatcher
func (di *Dispatcher) PreAuthentication(ctx context.Context, evt *CognitoPoolEventEnvelope) error {

	req := new(events.CognitoEventUserPoolsPreAuthenticationRequest)
	err := json.Unmarshal(*evt.Request, req)
	if err != nil {
		return err
	}

	res, err := di.ct.PreAuthentication(ctx, evt, req)
	if err != nil {
		return err
	}

	rawRes, err := RawResponse(res)
	if err != nil {
		return err
	}

	evt.Response = rawRes

	return nil
}

// PreToken pre token trigger dispatcher
func (di *Dispatcher) PreToken(ctx context.Context, evt *CognitoPoolEventEnvelope) error {

	req := new(events.CognitoEventUserPoolsPreTokenGenRequest)
	err := json.Unmarshal(*evt.Request, req)
	if err != nil {
		return err
	}

	res, err := di.ct.PreToken(ctx, evt, req)
	if err != nil {
		return err
	}

	rawRes, err := RawResponse(res)
	if err != nil {
		return err
	}

	evt.Response = rawRes

	return nil
}

// PostConfirmationSignUp Post confirmation signup trigger dispatcher
func (di *Dispatcher) PostConfirmationSignUp(ctx context.Context, evt *CognitoPoolEventEnvelope) error {

	req := new(events.CognitoEventUserPoolsPostConfirmationRequest)
	err := json.Unmarshal(*evt.Request, req)
	if err != nil {
		return err
	}

	res, err := di.ct.PostConfirmationSignUp(ctx, evt, req)
	if err != nil {
		return err
	}

	rawRes, err := RawResponse(res)
	if err != nil {
		return err
	}

	evt.Response = rawRes

	return nil
}

func RawResponse(res interface{}) (*json.RawMessage, error) {
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	j := new(json.RawMessage)
	err = j.UnmarshalJSON(data)

	return j, err
}
