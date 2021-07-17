module github.com/kubesys/kubernetes-scheduler

go 1.15

require (
	github.com/NVIDIA/gpu-monitoring-tools v0.0.0-20210412222843-d2e8de5a7ca2
	github.com/google/uuid v1.1.2
	github.com/kubesys/kubernetes-client-go v0.7.0
	github.com/sirupsen/logrus v1.8.1
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
)

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/kubernetes => k8s.io/kubernetes v0.20.2
)
