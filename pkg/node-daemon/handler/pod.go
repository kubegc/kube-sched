package handler

import (
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"kubesys.io/dl-scheduler/pkg/util"
	"sync"
)

type PodHandler struct {
	client *kubesys.KubernetesClient
	runningPods *util.SortedSet
	mu sync.Mutex
}


func NewPodHandler(client *kubesys.KubernetesClient) *PodHandler {
	return &PodHandler {
		client: client,
		runningPods: util.NewSortedSet(),
	}
}


func (ph *PodHandler) DoAdded(obj map[string]interface{}) {
	pods, err := ph.client.ListResources("Pod", "default")
	if err != nil {
		log.Errorf("error list pods, %s", err)
	}
	var pl v1.PodList
	err = pods.Into(&pl)
	runningPods := util.NewSortedSet()

	fmt.Println(pl)
	for _, pod := range pl.Items {
		for _, ref := range pod.OwnerReferences {
			fmt.Println(ref.Kind)
			if ref.Kind == "Task" {
				if pod.Status.Phase == v1.PodRunning {
					runningPods.Add(pod.Namespace + "/" + pod.Name)
				}
			}
		}
	}
	fmt.Println(runningPods)
	ph.mu.Lock()
	defer ph.mu.Unlock()

	for _, pod := range ph.runningPods.SortedKeys() {
		if !runningPods.Contains(pod) {
			ph.runningPods.Add(pod)
		} else {
			ph.runningPods.Delete(pod)
		}
	}
	fmt.Println(ph.runningPods)
}

func (ph *PodHandler) DoModified(obj map[string]interface{}) {
}

func (ph *PodHandler) DoDeleted(obj map[string]interface{}) {
}



func (ph *PodHandler) SyncFile() {

}



