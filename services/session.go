package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

func GetAwsSession() *session.Session {
	return session.Must(session.NewSession(
		&aws.Config{
			Endpoint: aws.String(os.Getenv("AWS_ENDPOINT")),
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
	))
}