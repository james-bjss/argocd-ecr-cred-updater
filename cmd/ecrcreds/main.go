package main

import (
	"ecrcredrotation/internal/config"
	"ecrcredrotation/internal/ecr"
	"ecrcredrotation/internal/k8s"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Parse Flags
	conf, output, err := config.ParseFlags(os.Args[0], os.Args[1:])

	// print default help
	if err == flag.ErrHelp {
		log.Info().Msg(output)
		os.Exit(2)
	}

	if err != nil {
		log.Error().Msgf("Error Parsing config line flags: %v", err)
		log.Error().Msg(output)
		os.Exit(1)
	}

	// Check Inputs
	validationErrs := conf.Validate()
	if len(validationErrs) > 0 {
		for _, e := range validationErrs {
			fmt.Println(e)
		}
		os.Exit(1)
	}

	log.Info().Msg("Fetching ECR Credentials...")
	log.Info().Msgf("Namespace: %s", conf.Namespace)
	log.Info().Msgf("Secret Name: %s", conf.SecretName)

	ecrClient, err := ecr.InitializeConfig(conf.AwsRegion)
	if err != nil {
		log.Error().Msg("Error Initialising ECR Client")
		os.Exit(1)
	}

	log.Info().Msgf("Fetching ECR Auth Token..")
	ecr := ecr.New(ecrClient)
	secret, err := ecr.GetPassword()

	if err != nil {
		log.Error().Msgf("Failed to Fetch ECR Token: %v", err)
		os.Exit(1)
	}

	var clientSet *kubernetes.Clientset
	if _, err := os.Stat(conf.KubeConfig); err == nil {
		clientSet, err = k8s.ConfigureClientWithConfig(&conf.KubeConfig)
	} else {
		clientSet, err = k8s.ConfigureInClusterClient()
	}

	k8sClient := k8s.New(clientSet, conf.Namespace)

	log.Info().Msgf("Patching Secret..")
	err = k8sClient.PatchSecret(conf.SecretName, secret)
	if err != nil {
		log.Error().Msgf("Failed to Patch Secret: %v", err)
		os.Exit(1)
	}
	log.Info().Msgf("Complete")
}
