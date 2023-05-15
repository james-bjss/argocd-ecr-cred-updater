package k8s

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	var client = testclient.NewSimpleClientset()
	New(client, "argocd")
}
func TestPatch(t *testing.T) {
	nsmock := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "argocd",
			CreationTimestamp: metav1.Date(2023, time.November, 10, 23, 0, 0, 0, time.UTC),
		},
	}

	secretmock := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test-secret",
			Namespace:         "argocd",
			CreationTimestamp: metav1.Date(2023, time.May, 8, 23, 1, 2, 0, time.UTC),
		},
		Data: map[string][]byte{
			"password": []byte("OldPassword"),
		},
	}

	var client = testclient.NewSimpleClientset(nsmock, secretmock)

	k8client := New(client, "argocd")
	err := k8client.PatchSecret("test-secret", "testing123")
	assert.NoError(t, err)

	result, _ := client.CoreV1().Secrets("argocd").Get(context.TODO(), "test-secret", metav1.GetOptions{})
	fmt.Println(string(result.Data["password"]))

	assert.Equal(t, "testing123", string(result.Data["password"]), "The password should be updated")
}
