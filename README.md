# ArgoCD ECR Credential Updater
Simple Utility to Patch ECR Credentials Used By ArgoCD Helm Sources. The tool fetches ECR credentials and patches the password field of the specified secret.
Out of the box ArgoCD does not rotate these See: [Issue #8096](https://github.com/argoproj/argo-cd/issues/8097)

## Local Usage
```console
# Passing AWS Region via command line flag
james@box:~$ ecrcredrotation -namespace argocd -secret some-secret --region us-east-1

# Passing AWS region by envar
james@box:~$ AWS_REGION=us-east-1; ecrcredrotation -namespace argocd -secret some-secret --region us-east-1

# Pass All params by envar
james@box:~$ AWS_REGION=us-east-1; ARGOCD_NAMESPACE=argocd; ARGOCD_ECR_SECRET=test-secret; ecrcredrotation
```

## Authentication

### Local 
The binary uses the default AWS search path for credentials.
If running locally the config assumes your default Kube credentials

### In-Cluster
The binary should be executed with a ServiceAccount which has write access to secrets in the target namespace and must have access to generate read credentials for the Target ECR repository.
