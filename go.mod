module github.com/kubesys/kubernetes-scheduler

go 1.15

require (
	github.com/kr/pretty v0.2.0 // indirect
	github.com/kubesys/kubernetes-client-go v1.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/sys v0.0.0-20210324051608-47abb6519492 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace k8s.io/kubernetes => k8s.io/kubernetes v0.20.2
