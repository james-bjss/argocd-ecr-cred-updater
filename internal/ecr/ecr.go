package ecr

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

type client struct {
	ecrClient ecriface.ECRAPI
}

func New(ecr ecriface.ECRAPI) *client {
	return &client{
		ecrClient: ecr,
	}
}

func InitializeConfig(region string) (ecriface.ECRAPI, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(region),
		CredentialsChainVerboseErrors: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}
	return ecr.New(sess), nil
}

func (c client) GetPassword() (string, error) {
	input := &ecr.GetAuthorizationTokenInput{}

	result, err := c.ecrClient.GetAuthorizationToken(input)

	if err != nil {
		return "", err
	}

	if len(result.AuthorizationData) == 0 {
		return "", errors.New("no ECR authorisation token returned")
	}
	return *result.AuthorizationData[0].AuthorizationToken, nil
}
