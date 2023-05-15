package cmd

import (
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseFlagsSuccess(t *testing.T) {
	var tests = []struct {
		name string
		args []string
		conf Config
	}{
		{"Flags All Args",
			[]string{"-namespace", "argocd", "-secret", "some-secret", "-kubeconfig", "/tmp/config", "-region", "us-east-1"},
			Config{
				Namespace:  "argocd",
				SecretName: "some-secret",
				KubeConfig: "/tmp/config",
				AwsRegion:  "us-east-1",
				Args:       []string{},
			},
		},
		{"Flags Default Kube Config Path",
			[]string{"-namespace", "argocd", "-secret", "some-secret"},
			Config{
				Namespace:  "argocd",
				SecretName: "some-secret",
				AwsRegion:  "",
				KubeConfig: filepath.Join(homedir.HomeDir(), ".kube", "config"),
				Args:       []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, output, err := ParseFlags("ecrcreds", tt.args)
			if err != nil {
				t.Errorf("err got %v, want nil", err)
			}
			if output != "" {
				t.Errorf("output got %q, want empty", output)
			}
			if !reflect.DeepEqual(*conf, tt.conf) {
				t.Errorf("conf got %+v, want %+v", *conf, tt.conf)
			}
		})
	}
}

func TestParseFlagsFailure(t *testing.T) {
	var tests = []struct {
		name string
		args []string
		conf Config
	}{
		{"No Args",
			[]string{},
			Config{
				Namespace:  "argocd",
				SecretName: "some-secret",
				KubeConfig: "/tmp/config",
				Args:       []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, output, err := ParseFlags("ecrcreds", tt.args)
			if err != nil {
				t.Errorf("err got %v, want nil", err)
			}
			if output != "" {
				t.Errorf("output got %q, want empty", output)
			}
		})
	}
}
