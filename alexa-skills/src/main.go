package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	alexa "github.com/ericdaugherty/alexa-skills-kit-golang"
	"github.com/ppco/smart-home/alexa-skills/src/switchbot"
)

var a = &alexa.Alexa{
	RequestHandler:      switchbot.NewHandler(),
	IgnoreApplicationID: true,
	IgnoreTimestamp:     true,
}

func Handle(ctx context.Context, requestEnv *alexa.RequestEnvelope) (any, error) {
	return a.ProcessRequest(ctx, requestEnv)
}

func main() {
	lambda.Start(Handle)
}
