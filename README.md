# kube-sched

## Running

Switch to the project root directory.

### Make kube-sched local binary
```
make
```

### Make kube-sched docker image
```
docker build -t doslab/kube-sched:v0.1-amd64 .
```

### Run kube-sched Pod
```
kubectl apply -f ./deploy/kube-sched.yaml
```