package ecr

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"testing"
	"time"
)

type mockECRClient struct {
	ecriface.ECRAPI
}

func String(v string) *string { return &v }

func (m *mockECRClient) GetAuthorizationToken(*ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
	expiry := time.Now().Add(time.Hour * 12)
	fmt.Println("called")
	return &ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			&ecr.AuthorizationData{
				AuthorizationToken: String("blaaah:blaaah"),
				ExpiresAt:          &expiry,
				ProxyEndpoint:      String("https://0000000000.dkr.ecr.us-east-1.amazonaws.com"),
			},
		},
	}, nil
}

func TestNew(t *testing.T) {
	mockSvc := &mockECRClient{}
	client := New(mockSvc)
	fmt.Println(client.GetPassword())
}
