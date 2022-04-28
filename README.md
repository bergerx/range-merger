
# range-merger

Accepts JSON array of range objects and return JSON array of "merged" range objects.

Example for payloads and returns (keep in mind, that there is no guarantee order of the input ranges).

`POST /api/overlaps`

```json
[{ "start": 3, "end": 5 }, { "start": 9, "end": 15 }, { "start": 4, "end": 7 }, 
 { "start": 10, "end": 12 }, { "start": 21, "end": 22 }, { "start": 17, "end": 21 }]
```

`output`

```json
[{"start": 3, "end": 7}, { "start": 9, "end": 15 }, { "start": 17, "end": 22 }]
```

## Some highlights

* The container image based on distroless, and its size is minimal (~25mb).
* App runs as nonroot in the conatiner.
* Example kubernetes deployment manifest comes with:
  * Runs as nonroot and `capabilities` are limited to `NET_BIND_SERVICE`.
  * liveness and readiness probes to help .
  * 3 instances deployed with `rollingUpdate` strategy with max surge 1.
  * `podAntiAffinity` rules to get the instances distributed to nodes whenever possible (using soft affinity),
    this can be improved by using `topologySpreadConstraints` in environment with multiple failure domains.
  * a `PodDisruptionBudget` that would help keeping HA during planned node operations.
## Build and run using local go

```shell
go build
```

## Build and run using docker

```shell
docker build -t range-merger:0.0.1 .
docker run --rm -p 8090:8090 range-merger:0.0.1
```

## Test a local environment

```shell
curl localhost:8090/api/overlaps -d '[{"start": 1, "end":3}, {"start": 2, "end": 4}]';
```

## Setup a local minikube environment

tested with minikube v1.25.2

```shell
docker build -t range-merger:0.0.1 .
minikube start -p range-merger --addons ingress 
minikube image -p range-merger load range-merger:0.0.1
kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s
kubectl apply -f deploy.yaml
kubectl wait --for=condition=ready pod --selector=app.kubernetes.io/name=range-merger --timeout=90s
sleep 5  # ingress controller takes some time to reconfigure after the pod becomes ready 
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 8090:80 &  # this is a hack for minikube ingress
curl range-merger.127.0.0.1.nip.io:8090/api/overlaps -d '[{ "start": 3, "end": 5 }, { "start": 9, "end": 15 }, { "start": 4, "end": 7 }, { "start": 10, "end": 12 }, { "start": 21, "end": 22 }, { "start": 17, "end": 21 }]'

kill %1
# delete the test environment once you are done
minikube delete -p range-merger
```

## Setup a local kind environment

tested with kind v0.12.0

```shell
docker build -t range-merger:0.0.1 .
kind create cluster --name range-merger --config kind.config
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml 
kind load docker-image range-merger:0.0.1 --name range-merger
kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s
kubectl apply -f deploy.yaml
kubectl wait --for=condition=ready pod --selector=app.kubernetes.io/name=range-merger --timeout=90s
sleep 5  # ingress controller takes some time to reconfigure after the pod becomes ready 
curl range-merger.127.0.0.1.nip.io/api/overlaps -d '[{ "start": 3, "end": 5 }, { "start": 9, "end": 15 }, { "start": 4, "end": 7 }, { "start": 10, "end": 12 }, { "start": 21, "end": 22 }, { "start": 17, "end": 21 }]'

# delete the test environment once you are done
kind delete cluster --name range-merger 
```