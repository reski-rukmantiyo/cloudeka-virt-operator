# Development

## Development Environment

- Kind 
- Docker
- Helm
- Kubectl

## Development Process

```
operator-sdk init --domain=cloudeka.ai --repo=github.com/reski-rukmantiyo/cloudeka-virt-operator
operator-sdk create api --group cloudeka.ai --version v1alpha1 --kind CloudekaMachine --resource --controller
operator-sdk create webhook --group cloudeka.ai --version v1alpha1 --kind CloudekaMachine --defaulting --programmatic-validation
```

## Development Tools

Visual Studio Code

## Generate CRD

```
make generate manifests install
```