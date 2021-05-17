package handler

import (
	"bufio"
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"kubesys.io/dl-scheduler/pkg/util"
	"os"
	"strconv"
	"sync"
)

type PodHandler struct {
	client *kubesys.KubernetesClient
	//runningPods util.SortedSet
	gpu2RunningPodList map[string]*util.SortedSet
	mu sync.Mutex
}


func NewPodHandler(client *kubesys.KubernetesClient) *PodHandler {
	return &PodHandler {
		client: client,
		gpu2RunningPodList: make(map[string]*util.SortedSet),
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
	//jb, _ := json.Marshal(obj)

	pods, err := ph.client.ListResources("Pod", "default")
	if err != nil {
		log.Errorf("error list pods, %s", err)
	}
	var pl v1.PodList
	err = pods.Into(&pl)
	gpu2RunningPods := make(map[string][]string)

	//fmt.Println(pl)
	for _, pod := range pl.Items {
		for _, ref := range pod.OwnerReferences {
			if ref.Kind == "Task" {
				if pod.Status.Phase == v1.PodRunning {
					uuid := pod.Annotations["schedule-gpuid"]
					if gpu2RunningPods[uuid] == nil {
						gpu2RunningPods[uuid] = make([]string, 0)
					}
					gpu2RunningPods[uuid] = append(gpu2RunningPods[uuid], pod.Namespace + "/" + pod.Name)
				}
			}
		}
	}
	fmt.Println(gpu2RunningPods)
	ph.mu.Lock()
	defer ph.mu.Unlock()
	for gpu, tasks := range gpu2RunningPods {
		lastTasks := ph.gpu2RunningPodList[gpu]
		if lastTasks == nil {
			lastTasks = util.NewSortedSet()
		}
		if !util.Compare(tasks, lastTasks.SortedKeys()) {
			ph.gpu2RunningPodList[gpu] = util.NewSortedSet()
			for _, task := range tasks {
				ph.gpu2RunningPodList[gpu].Add(task)
			}
			fmt.Println("sync called")
			ph.SyncFile()
		}
	}
	//
	//for _, pod := range ph.runningPods.SortedKeys() {
	//	if !runningPods.Contains(pod) {
	//		ph.runningPods.Delete(pod)
	//	}
	//}
	//
	//for _, pod := range runningPods.SortedKeys() {
	//	uuid := pod
	//	if !ph.runningPods.Contains(pod) {
	//		ph.runningPods.Add(pod)
	//	}
	//}
	//fmt.Println("Running:", ph.runningPods)


}

func (ph *PodHandler) SyncFile() {
	for gpu, tasks := range ph.gpu2RunningPodList {
		configFile, err := os.Create("/Users/yangchen/kubeshare/scheduler/config/" + gpu)
		if err != nil {
			fmt.Println("file create error:", err)
		}
		buf := bufio.NewWriter(configFile)
		buf.Write([]byte(strconv.Itoa(tasks.Size())))
		buf.Write([]byte("\n"))
		for _, task := range tasks.SortedKeys() {
			buf.Write([]byte(task))
			// TODO: write control info
			buf.Write([]byte("\n"))
		}
		_ = buf.Flush()
		configFile.Sync()
		configFile.Close()
	}

	for gpu, tasks := range ph.gpu2RunningPodList {
		portFile, err := os.Create("/Users/yangchen/kubeshare/scheduler/podmanagerport/" + gpu)
		if err == nil {
			fmt.Println("file create error:", err)
		}
		buf := bufio.NewWriter(portFile)
		buf.Write([]byte(strconv.Itoa(tasks.Size())))
		buf.Write([]byte("\n"))
		for _, task := range tasks.SortedKeys() {
			buf.Write([]byte(task))
			// TODO: write port info
			buf.Write([]byte("\n"))
		}
		_ = buf.Flush()
		portFile.Sync()
		portFile.Close()
	}
}



