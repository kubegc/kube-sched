package handler
//
//import (
//	"bufio"
//	"fmt"
//	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
//	log "github.com/sirupsen/logrus"
//	v1 "k8s.io/api/core/v1"
//	"kubesys.io/dl-scheduler/pkg/util"
//	"os"
//	"strconv"
//	"sync"
//)
//
//type PodHandler struct {
//	client *kubesys.KubernetesClient
//	//runningPods util.SortedSet
//	gpu2RunningPodList map[string]*util.SortedSet
//	gpuPod2Request map[string]map[string]string
//	gpuPod2Limit map[string]map[string]string
//	gpuPod2Memory map[string]map[string]string
//	gpuPod2Port map[string]map[string]string
//	mu sync.Mutex
//}
//
//
//func NewPodHandler(client *kubesys.KubernetesClient) *PodHandler {
//	return &PodHandler {
//		client: client,
//		gpu2RunningPodList: make(map[string]*util.SortedSet),
//		gpuPod2Request: make(map[string]map[string]string),
//		gpuPod2Limit: make(map[string]map[string]string),
//		gpuPod2Memory: make(map[string]map[string]string),
//	}
//}
//
//
//func (ph *PodHandler) DoAdded(obj map[string]interface{}) {
//	ph.syncHandler(obj)
//}
//
//func (ph *PodHandler) DoModified(obj map[string]interface{}) {
//	ph.syncHandler(obj)
//}
//
//func (ph *PodHandler) DoDeleted(obj map[string]interface{}) {
//	ph.syncHandler(obj)
//}
//
//func (ph *PodHandler) syncHandler(obj map[string]interface{}) {
//	//jb, _ := json.Marshal(obj)
//
//	pods, err := ph.client.ListResources("Pod", "default")
//	if err != nil {
//		log.Errorf("error list pods, %s", err)
//	}
//	var pl v1.PodList
//	err = pods.Into(&pl)
//	gpu2RunningPods := make(map[string][]string)
//	gpuPod2Request := make(map[string]map[string]string)
//	gpuPod2Limit := make(map[string]map[string]string)
//	gpuPod2Memory:= make(map[string]map[string]string)
//
//	//fmt.Println(pl)
//	for _, pod := range pl.Items {
//		for _, ref := range pod.OwnerReferences {
//			if ref.Kind == "Task" {
//				if pod.Status.Phase == v1.PodRunning {
//					uuid := pod.Annotations["schedule-gpuid"]
//					req := pod.Annotations["gpu-request"]
//					limit := pod.Annotations["gpu-limit"]
//					memory := pod.Annotations["gpu-memory"]
//
//					podName := pod.Namespace + "/" + pod.Name
//					if gpu2RunningPods[uuid] == nil {
//						gpu2RunningPods[uuid] = make([]string, 0)
//					}
//					gpu2RunningPods[uuid] = append(gpu2RunningPods[uuid], podName)
//
//					if gpuPod2Request[uuid] == nil {
//						gpuPod2Request[uuid] = make(map[string]string)
//					}
//					gpuPod2Request[uuid][podName] = req
//
//					if gpuPod2Limit[uuid] == nil {
//						gpuPod2Limit[uuid] = make(map[string]string)
//					}
//					gpuPod2Limit[uuid][podName] = limit
//
//					if gpuPod2Memory[uuid] == nil {
//						gpuPod2Memory[uuid] = make(map[string]string)
//					}
//					gpuPod2Memory[uuid][podName] = memory
//				}
//			}
//		}
//	}
//	fmt.Println(gpu2RunningPods)
//	ph.mu.Lock()
//	defer ph.mu.Unlock()
//	for gpu, tasks := range gpu2RunningPods {
//		lastTasks := ph.gpu2RunningPodList[gpu]
//		if lastTasks == nil {
//			lastTasks = util.NewSortedSet()
//		}
//		if !util.Compare(tasks, lastTasks.SortedKeys()) {
//			ph.gpu2RunningPodList[gpu] = util.NewSortedSet()
//			for _, task := range tasks {
//				ph.gpu2RunningPodList[gpu].Add(task)
//			}
//			fmt.Println("sync called")
//			ph.SyncFile(gpuPod2Request, gpuPod2Limit, gpuPod2Memory)
//		}
//	}
//	//
//	//for _, pod := range ph.runningPods.SortedKeys() {
//	//	if !runningPods.Contains(pod) {
//	//		ph.runningPods.Delete(pod)
//	//	}
//	//}
//	//
//	//for _, pod := range runningPods.SortedKeys() {
//	//	uuid := pod
//	//	if !ph.runningPods.Contains(pod) {
//	//		ph.runningPods.Add(pod)
//	//	}
//	//}
//	//fmt.Println("Running:", ph.runningPods)
//
//
//}
//
//func (ph *PodHandler) SyncFile(req, limit, memory map[string]map[string]string) {
//	for gpu, tasks := range ph.gpu2RunningPodList {
//		configFile, err := os.Create("/Users/yangchen/kubeshare/scheduler/config/" + gpu)
//		if err != nil {
//			fmt.Println("file create error:", err)
//		}
//		buf := bufio.NewWriter(configFile)
//		buf.Write([]byte(strconv.Itoa(tasks.Size())))
//		buf.Write([]byte("\n"))
//		for _, task := range tasks.SortedKeys() {
//			buf.Write([]byte(task))
//			buf.Write([]byte(" "))
//			buf.Write([]byte(req[gpu][task]))
//			buf.Write([]byte(" "))
//			buf.Write([]byte(limit[gpu][task]))
//			buf.Write([]byte(" "))
//			buf.Write([]byte(memory[gpu][task]))
//			buf.Write([]byte("\n"))
//		}
//		_ = buf.Flush()
//		configFile.Sync()
//		configFile.Close()
//	}
//
//	for gpu, tasks := range ph.gpu2RunningPodList {
//		portFile, err := os.Create("/Users/yangchen/kubeshare/scheduler/podmanagerport/" + gpu)
//		if err == nil {
//			fmt.Println("file create error:", err)
//		}
//		buf := bufio.NewWriter(portFile)
//		buf.Write([]byte(strconv.Itoa(tasks.Size())))
//		buf.Write([]byte("\n"))
//		for _, task := range tasks.SortedKeys() {
//			buf.Write([]byte(task))
//			// TODO: write port info
//
//			buf.Write([]byte("\n"))
//		}
//		_ = buf.Flush()
//		portFile.Sync()
//		portFile.Close()
//	}
//}
//
//
//
//
//
