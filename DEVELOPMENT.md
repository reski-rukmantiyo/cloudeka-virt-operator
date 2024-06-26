# Development

## Development Environment

- Kind 
- Docker
- Helm
- Kubectl

## Preqrequisite

### Create SSL Certificate (development only)

```
openssl genrsa -out tls.key 2048
openssl req -new -key tls.key -out tls.csr -subj "/CN=my-webhook.local" -addext "subjectAltName = DNS:my-webhook.local,DNS:localhost"
openssl x509 -req -in tls.csr -signkey tls.key -out tls.crt -days 365 -extfile <(printf "subjectAltName=DNS:my-webhook.local,DNS:localhost")
cp tls.crt /tmp/k8s-webhook-server/serving-certs/tls.crt
cp tls.key /tmp/k8s-webhook-server/serving-certs/tls.key
cat /tmp/k8s-webhook-server/serving-certs/tls.crt|base64|tr -d '\n'
```

### Setup Webhook Server (development only)

```
10.0.2.2 my-webhook.local
```

## Development Process

```
operator-sdk init --domain=cloudeka.ai --repo=github.com/reski-rukmantiyo/cloudeka-virt-operator
operator-sdk create api --group virt --version v1alpha1 --kind CloudekaMachine --resource --controller
operator-sdk create webhook --group virt --version v1alpha1 --kind CloudekaMachine --defaulting --programmatic-validation
```

## Development Tools

Visual Studio Code

## Generate CRD

```
make manifests generate install
```

