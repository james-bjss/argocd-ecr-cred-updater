package config

import (
	"bytes"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type Config struct {
	Namespace  string
	SecretName string
	KubeConfig string
	AwsRegion  string
	Args       []string
}

func (c Config) Validate() (errors []string) {
	// Check Inputs
	validationErrs := validation.IsDNS1123Label(c.Namespace)
	if len(validationErrs) > 0 {
		return append([]string{fmt.Sprintf("invalid namespace name: %s", c.Namespace)}, validationErrs...)
	}

	validationErrs = validation.IsDNS1123Subdomain(c.SecretName)
	if len(validationErrs) > 0 {
		return append([]string{fmt.Sprintf("invalid secret name: %s", c.SecretName)}, validationErrs...)
	}
	return nil
}

func ParseFlags(progname string, args []string) (config *Config, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var conf Config
	flags.StringVar(&conf.Namespace, "namespace", os.Getenv("ARGOCD_NAMESPACE"), "Namespace where ArgoCD ECR secret resides")
	flags.StringVar(&conf.SecretName, "secret", os.Getenv("ARGOCD_ECR_SECRET"), "Name of ArgoCD Secret to Patch")
	flags.StringVar(&conf.AwsRegion, "region", "", "AWS Region to use")
	flags.StringVar(&conf.KubeConfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	conf.Args = flags.Args()
	return &conf, buf.String(), nil
}
