package node_daemon

import (
	"encoding/json"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/workqueue"
	v1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"fmt"
	"kubesys.io/dl-scheduler/pkg/scheduler"
	"kubesys.io/dl-scheduler/pkg/util"
)

type Worker struct {
	client *kubesys.KubernetesClient
	workqueue workqueue.RateLimitingInterface
}

func NewWorker(client *kubesys.KubernetesClient, wq workqueue.RateLimitingInterface) *Worker{
	return &Worker{
		client:    client,
		workqueue: wq,
	}
}

func (w *Worker) Run() {
	for {
		obj, shutdown := w.workqueue.Get()
		fmt.Println(obj)
		if shutdown {
			log.Errorf("node-daemon worker shutdown")
			break
		}
		w.syncHandler(obj)
	}
}


func (w *Worker) syncHandler(obj interface{}) {
	var task v1.Task
	taskStr, _ := obj.(string)
	json.Unmarshal([]byte(taskStr), &task)

	pod, err := w.client.GetResource("Pod", task.Namespace, task.Name)
	if err != nil {
		log.Errorf("error get pod, %s", err)
	}

	ss := pod.GetString("reason")
	if ss == "NotFound" && util.Scheduled(&task){
		fmt.Println(pod)
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
				//corev1.EnvVar{
				//	Name:  "LD_PRELOAD",
				//	Value: KubeShareLibraryPath + "/libgemhook.so.1",
				//},
				//corev1.EnvVar{
				//	Name:  "POD_MANAGER_IP",
				//	Value: podManagerIP,
				//},
				//corev1.EnvVar{
				//	Name:  "POD_MANAGER_PORT",
				//	Value: fmt.Sprintf("%d", podManagerPort),
				//},
				corev1.EnvVar{
					Name:  "POD_NAME",
					Value: fmt.Sprintf("%s/%s", task.ObjectMeta.Namespace, task.ObjectMeta.Name),
				},
			)
			//c.VolumeMounts = append(c.VolumeMounts,
			//	corev1.VolumeMount{
			//		Name:      "kubeshare-lib",
			//		MountPath: KubeShareLibraryPath,
			//	},
			//)
		}
		//specCopy.Volumes = append(specCopy.Volumes,
		//	corev1.Volume{
		//		Name: "kubeshare-lib",
		//		VolumeSource: corev1.VolumeSource{
		//			HostPath: &corev1.HostPathVolumeSource{
		//				Path: KubeShareLibraryPath,
		//			},
		//		},
		//	},
		//)
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      task.ObjectMeta.Name,
				Namespace: task.ObjectMeta.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(&task, schema.GroupVersionKind{
						Group:   v1.GroupVersion.Group,
						Version: v1.GroupVersion.Version,
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
		_, err = w.client.CreateResource(string(podByte))
		if err != nil {
			log.Errorf("error create new pod, %s", err)
		}
	}
}
