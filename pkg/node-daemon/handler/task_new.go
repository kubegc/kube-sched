package handler

import (
	"bufio"
	"encoding/json"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dosv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	node_daemon "kubesys.io/dl-scheduler/pkg/node-daemon"
	"kubesys.io/dl-scheduler/pkg/scheduler"
	"kubesys.io/dl-scheduler/pkg/util"
	"os"
	"strconv"
	"sync"
	log "github.com/sirupsen/logrus"
	"fmt"
)

type TaskHandler struct {
	client *kubesys.KubernetesClient
	gpu2RunningTasks map[string]*util.SortedSet
	gpuPod2Request map[string]map[string]string
	gpuPod2Limit map[string]map[string]string
	gpuPod2Memory map[string]map[string]string
	task2Port map[string]int
	basePort int
	freePort *util.Bitmap64
	mu sync.Mutex
}

func NewTaskHandler(client *kubesys.KubernetesClient) *TaskHandler {
	return &TaskHandler{
		client: client,
		gpu2RunningTasks: make(map[string]*util.SortedSet),
		gpuPod2Request: make(map[string]map[string]string),
		gpuPod2Limit: make(map[string]map[string]string),
		gpuPod2Memory: make(map[string]map[string]string),
		task2Port: make(map[string]int),
		basePort: 50051,
		freePort: util.NewBitMap64(1000),
	}
}

func (th *TaskHandler) DoAdded(obj map[string]interface{}) {
	th.syncHandler(obj)
}

func (th *TaskHandler) DoModified(obj map[string]interface{}) {
	th.syncHandler(obj)
}

func (th *TaskHandler) DoDeleted(obj map[string]interface{}) {
	th.syncHandler(obj)
}



func (th *TaskHandler) syncHandler(obj map[string]interface{}) {

	tasks, err := th.client.ListResources("Task", "default")
	if err != nil {
		log.Errorf("error list pods, %s", err)
	}
	var tl dosv1.TaskList
	err = tasks.Into(&tl)
	th.mu.Lock()
	allTask := make(map[string]bool)

	for _, task := range tl.Items {
		taskName := task.Namespace + "/" + task.Name
		allTask[taskName] = true
	}

	for task, _ := range allTask {
		if _, ok := th.task2Port[task]; !ok {
			th.task2Port[task] = th.freePort.Acquire()
		}
	}

	for taskName, port := range th.task2Port {
		if _, ok := allTask[taskName]; !ok {
			th.freePort.Release(port)
			delete(th.task2Port, taskName)
		}
	}


	th.mu.Unlock()
	fmt.Println(th.task2Port)
	// Create corresponding pod if pod does not exist
	var task dosv1.Task
	jsonByte, _ := json.Marshal(obj)
	_ = json.Unmarshal(jsonByte, &task)
	pod, err := th.client.GetResource("Pod", task.Namespace, task.Name)
	if err != nil {
		log.Errorf("error get pod, %s", err)
	}
	if  pod == nil && util.Scheduled(&task) {
		th.CreatePod(&task)
	}
	// Sync running tasks
	gpu2RunningTasks := make(map[string]*util.SortedSet)
	gpuPod2Request := make(map[string]map[string]string)
	gpuPod2Limit := make(map[string]map[string]string)
	gpuPod2Memory:= make(map[string]map[string]string)
	//gpuPod2Port := make(map[string]map[string]string)

	//fmt.Println(pl)
	for _, task := range tl.Items {
		if task.Status.PodStatus.Phase == corev1.PodRunning {
			uuid := task.Annotations["schedule-gpuid"]
			req := task.Annotations["gpu-request"]
			limit := task.Annotations["gpu-limit"]
			memory := task.Annotations["gpu-memory"]
			taskName := task.Namespace + "/" + task.Name
			if gpu2RunningTasks[uuid] == nil {
				gpu2RunningTasks[uuid] = util.NewSortedSet()
			}
			gpu2RunningTasks[uuid].Add(taskName)

			if gpuPod2Request[uuid] == nil {
				gpuPod2Request[uuid] = make(map[string]string)
			}
			gpuPod2Request[uuid][taskName] = req

			if gpuPod2Limit[uuid] == nil {
				gpuPod2Limit[uuid] = make(map[string]string)
			}
			gpuPod2Limit[uuid][taskName] = limit

			if gpuPod2Memory[uuid] == nil {
				gpuPod2Memory[uuid] = make(map[string]string)
			}
			gpuPod2Memory[uuid][taskName] = memory

			//if gpuPod2Port[uuid] == nil {
			//	gpuPod2Port[uuid] = make(map[string]string)
			//}
			//gpuPod2Port[uuid][taskName] =
		}
		//for _, ref := range task.OwnerReferences {
		//	if ref.Kind == "Task" {
		//		if pod.Status.Phase == v1.PodRunning {
		//			uuid := pod.Annotations["schedule-gpuid"]
		//			req := pod.Annotations["gpu-request"]
		//			limit := pod.Annotations["gpu-limit"]
		//			memory := pod.Annotations["gpu-memory"]
		//
		//			podName := pod.Namespace + "/" + pod.Name
		//			if gpu2RunningPods[uuid] == nil {
		//				gpu2RunningPods[uuid] = make([]string, 0)
		//			}
		//			gpu2RunningPods[uuid] = append(gpu2RunningPods[uuid], podName)
		//
		//			if gpuPod2Request[uuid] == nil {
		//				gpuPod2Request[uuid] = make(map[string]string)
		//			}
		//			gpuPod2Request[uuid][podName] = req
		//
		//			if gpuPod2Limit[uuid] == nil {
		//				gpuPod2Limit[uuid] = make(map[string]string)
		//			}
		//			gpuPod2Limit[uuid][podName] = limit
		//
		//			if gpuPod2Memory[uuid] == nil {
		//				gpuPod2Memory[uuid] = make(map[string]string)
		//			}
		//			gpuPod2Memory[uuid][podName] = memory
		//		}
		//	}
		//}
	}
	fmt.Println(gpu2RunningTasks)
	th.mu.Lock()
	defer th.mu.Unlock()
	for gpu, tasks := range gpu2RunningTasks {
		lastTasks := th.gpu2RunningTasks[gpu]
		if lastTasks == nil {
			lastTasks = util.NewSortedSet()
		}
		if !util.Compare(tasks.SortedKeys(), lastTasks.SortedKeys()) {
			th.gpu2RunningTasks[gpu] = util.NewSortedSet()
			for _, task := range tasks.SortedKeys() {
				th.gpu2RunningTasks[gpu].Add(task)
			}


			fmt.Println("sync called")
			th.SyncFile(gpuPod2Request, gpuPod2Limit, gpuPod2Memory)
		}
	}
	//

	for gpu, tasks := range th.gpu2RunningTasks {
		for _, task := range tasks.SortedKeys() {
			if !gpu2RunningTasks[gpu].Contains(task) {
				th.gpu2RunningTasks[gpu].Delete(task)

			}
		}
	}

	for gpu, tasks := range gpu2RunningTasks {
		for _, task := range tasks.SortedKeys() {
			if !th.gpu2RunningTasks[gpu].Contains(task) {
				th.gpu2RunningTasks[gpu].Add(task)
			}
		}

	}
	//for _, pod := range th.runningPods.SortedKeys() {
	//	if !runningPods.Contains(pod) {
	//		th.runningPods.Delete(pod)
	//	}
	//}

	//for _, pod := range runningPods.SortedKeys() {
	//	uuid := pod
	//	if !th.runningPods.Contains(pod) {
	//		th.runningPods.Add(pod)
	//	}
	//}
	//fmt.Println("Running:")
	//for gpu, tasks := range th.gpu2RunningTasks {
	//	fmt.Println(gpu, tasks.SortedKeys())
	//}






}

func (th *TaskHandler) SyncFile(req, limit, memory map[string]map[string]string) {
	for gpu, tasks := range th.gpu2RunningTasks {
		configFile, err := os.Create("/Users/yangchen/kubeshare/scheduler/config/" + gpu)
		if err != nil {
			fmt.Println("file create error:", err)
		}
		buf := bufio.NewWriter(configFile)
		buf.Write([]byte(strconv.Itoa(tasks.Size())))
		buf.Write([]byte("\n"))
		for _, task := range tasks.SortedKeys() {
			buf.Write([]byte(task))
			buf.Write([]byte(" "))
			buf.Write([]byte(req[gpu][task]))
			buf.Write([]byte(" "))
			buf.Write([]byte(limit[gpu][task]))
			buf.Write([]byte(" "))
			buf.Write([]byte(memory[gpu][task]))
			buf.Write([]byte("\n"))
		}
		_ = buf.Flush()
		configFile.Sync()
		configFile.Close()
	}

	for gpu, tasks := range th.gpu2RunningTasks {
		portFile, err := os.Create("/Users/yangchen/kubeshare/scheduler/podmanagerport/" + gpu)
		if err == nil {
			fmt.Println("file create error:", err)
		}
		buf := bufio.NewWriter(portFile)
		buf.Write([]byte(strconv.Itoa(tasks.Size())))
		buf.Write([]byte("\n"))
		for _, task := range tasks.SortedKeys() {
			buf.Write([]byte(task))
			buf.Write([]byte(" "))
			buf.Write([]byte(strconv.Itoa(th.task2Port[task] + th.basePort)))
			buf.Write([]byte("\n"))
		}
		_ = buf.Flush()
		portFile.Sync()
		portFile.Close()
	}
}



func (th *TaskHandler) CreatePod(task *dosv1.Task) {
	specCopy := task.Spec.DeepCopy()
	annotationCopy := make(map[string]string, len(task.ObjectMeta.Annotations) + 4)
	for key, val := range task.ObjectMeta.Annotations {
		annotationCopy[key] = val
	}
	for i := range specCopy.Containers {
		c := &specCopy.Containers[i]
		c.Env = append(c.Env,
			corev1.EnvVar{
				Name:  "NVIDIA_VISIBLE_DEVICES",
				Value: task.Annotations[scheduler.ScheduleGPUIDAnnotation],
			},
			corev1.EnvVar{
				Name:  "NVIDIA_DRIVER_CAPABILITIES",
				Value: "compute,utility",
			},
			corev1.EnvVar{
				Name:  "LD_PRELOAD",
				Value: node_daemon.GemLibraryPath+ "/libgemhook.so.1",
			},
			corev1.EnvVar{
				Name:  "POD_MANAGER_IP",
				Value: "192.168.228.108",
			},
			corev1.EnvVar{
				Name:  "POD_MANAGER_PORT",
				Value: fmt.Sprintf("%d", 55051),
			},
			corev1.EnvVar{
				Name:  "POD_NAME",
				Value: fmt.Sprintf("%s/%s", task.ObjectMeta.Namespace, task.ObjectMeta.Name),
			},
		)
		c.VolumeMounts = append(c.VolumeMounts,
			corev1.VolumeMount{
				Name:      "kubeshare-lib",
				MountPath: node_daemon.GemLibraryPath,
			},
		)
	}
	specCopy.Volumes = append(specCopy.Volumes,
		corev1.Volume{
			Name: "kubeshare-lib",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: node_daemon.GemLibraryPath,
				},
			},
		},
	)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      task.ObjectMeta.Name,
			Namespace: task.ObjectMeta.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(task, schema.GroupVersionKind{
					Group:   dosv1.GroupVersion.Group,
					Version: dosv1.GroupVersion.Version,
					Kind:    "Task",
				}),
			},
			Annotations: annotationCopy,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		Spec: *specCopy,
	}
	podByte, err := json.Marshal(pod)
	if err != nil {
		log.Errorf("error marshal pod %s", err)
	}
	_, err = th.client.CreateResource(string(podByte))
	if err != nil {
		log.Errorf("error create new pod, %s", err)
	}
}