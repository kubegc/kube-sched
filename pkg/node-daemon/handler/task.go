package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	"k8s.io/client-go/util/workqueue"
)

type TaskHandler struct {
	client *kubesys.KubernetesClient
	queue workqueue.RateLimitingInterface
}


func NewTaskHandler(client *kubesys.KubernetesClient) *TaskHandler {
	return &TaskHandler{
		client: client,
		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "tasks"),
	}
}


func (th *TaskHandler) DoAdded(obj map[string]interface{}) {
	jb, _ := json.Marshal(obj)
	fmt.Println(string(jb))
}

func (th *TaskHandler) DoModified(obj map[string]interface{}) {

}

func (th *TaskHandler) DoDeleted(obj map[string]interface{}) {
	jb, _ := json.Marshal(obj)
	fmt.Println(string(jb))
}