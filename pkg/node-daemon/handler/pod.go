package handler

import (
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
)

type PodHandler struct {
	c *kubesys.KubernetesClient
}


func NewPodHandler(client *kubesys.KubernetesClient) *PodHandler {
	return &PodHandler{c: client}
}


func (ph *PodHandler) DoAdded(obj map[string]interface{}) {
	fmt.Println(obj)
}

func (ph *PodHandler) DoModified(obj map[string]interface{}) {
	fmt.Println(obj)
}

func (ph *PodHandler) DoDeleted(obj map[string]interface{}) {
	fmt.Println(obj)
}



