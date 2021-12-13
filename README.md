# kube-sched

## Running

Switch to the project root directory.

### Make kube-sched local binary
```
make
```

### Make kube-sched docker image
```
docker build -t registry.cn-beijing.aliyuncs.com/doslab/kube-sched:v0.3.3-amd64 .
```

### Run kube-sched Pod
```
kubectl apply -f ./deploy/kube-sched.yaml
```