# NATS-simple-demo


## Install NATS

```bash
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm install -n nats nats nats/nats -f ./nats.values.yaml --create-namespace
```
