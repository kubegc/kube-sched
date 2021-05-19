package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	v1 "k8s.io/api/core/v1"
	dosv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
)

type PodHandler struct {
	client *kubesys.KubernetesClient
}


func NewPodHandler(client *kubesys.KubernetesClient) *PodHandler {
	return &PodHandler {
		client: client,
	}
}


func (ph *PodHandler) DoAdded(obj map[string]interface{}) {
	ph.syncHandler(obj)
}

func (ph *PodHandler) DoModified(obj map[string]interface{}) {
	ph.syncHandler(obj)
}

func (ph *PodHandler) DoDeleted(obj map[string]interface{}) {
	ph.syncHandler(obj)
}

func (ph *PodHandler) syncHandler(obj map[string]interface{}) {
	var pod v1.Pod
	jsonByte, _ := json.Marshal(obj)
	_ = json.Unmarshal(jsonByte, &pod)
	ownedByTask := false
	for _, ref := range pod.OwnerReferences {
		if ref.Kind == "Task" {
			ownedByTask = true
		}
	}
	if !ownedByTask {
		return
	}
	var task dosv1.Task
	taskObj, err := ph.client.GetResource("Task", "default", pod.Name)
	if err != nil {
		fmt.Println(err)
	}

	if taskObj == nil {
		fmt.Println(err)
	}
	taskObj.Into(&task)

	task.Status.PodStatus.Phase = pod.Status.Phase

	_, err = ph.client.UpdateResourceObject(&task)

	if err != nil {
		fmt.Println("error update resource", err)
	}

}
