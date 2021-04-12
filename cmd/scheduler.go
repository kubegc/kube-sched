package main
import (
	"context"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	"k8s.io/client-go/util/workqueue"
	"kubesys.io/dl-scheduler/pkg/scheduler"
)

var (
	masterURL = "https://133.133.135.42:6443"
	token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjN2U3VZUk16R3ZfZGNaMkw4bVktVGlRWnJGZFB2NWprU1lrd0hObnNBVFEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbi1sNXhnayIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjQzYzM1YTYxLTgzYWMtNGJiYy1iNDE5LTdjZmFjODA5NzAyNCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.C7wtJgeDZgotLzw2Jl5cKPLRHQse2jM3Fx3yH6O82tR-pvs5yQuutCJu65y_81OJlbJzge0kzy8r6NeMJz_mDWVa5pGW4SuWga5IDs17PolzYwHJClLTpfMn1uYkOVna5aDKv8r7NdwiWgmvZ0aA4Ll7SNmDiNTaNPXT8cs-LDB1fisQsNQHiSu8ucl0PZDFfhFbH3QOlvkLRsZAoXRnzfNzgAaFJ7svd07JtJOeLTG53YZiX1WfyqxFR-dYP9hCG250SuNqozPHWxz_Z6ITQXwbJ_f6mTtLZXldLGGSKnvEJxmpyI1QbSaRaTGrWFj8fF5cBKANhTBceDMnotM4eg"
)

func main() {
	client := kubesys.NewKubernetesClient(masterURL, token)
	client.Init()

	workqueue :=  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "SharePods")
	controller := scheduler.NewController(client, workqueue)
	watcher := kubesys.NewKubernetesWatcher(client, scheduler.NewWorkQueueHandler(workqueue))
	ctx := context.Background()
	go controller.Run()
	client.WatchResources("Pod", "default", watcher)
	<-ctx.Done()
}