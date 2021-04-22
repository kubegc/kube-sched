package handler

import (
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"

)

type TaskHandler struct {
	c *kubesys.KubernetesClient
}


func NewTaskHandler(client *kubesys.KubernetesClient) *TaskHandler {
	return &TaskHandler{c: client}
}


func (th *TaskHandler) DoAdded(obj map[string]interface{}) {
	fmt.Println(obj)
}

func (th *TaskHandler) DoModified(obj map[string]interface{}) {
	fmt.Println(obj)
}

func (th *TaskHandler) DoDeleted(obj map[string]interface{}) {
	fmt.Println(obj)
}