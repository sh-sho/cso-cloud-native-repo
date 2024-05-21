# workload identity

## how to exec

### local

```sh
export NAMESPACE=<your-namespace>
export COMPARTMENT_ID=<your-compartment-id>
go run main.go
```

### build and push

```sh
docker image build -t <region-key>.ocir.io/orasejapan/<prefix>/workload-identity:1.0.0 .; \
docker image push <region-key>.ocir.io/orasejapan/<prefix>/workload-identity:1.0.0
```
