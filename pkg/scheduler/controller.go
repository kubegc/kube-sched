package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	dosv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler/algorithm"
	"strconv"
	"time"
)

type Controller struct {
	client *kubesys.KubernetesClient
	workqueue workqueue.RateLimitingInterface
}

func NewController(client *kubesys.KubernetesClient, workqueue workqueue.RateLimitingInterface) *Controller {
	return &Controller{
		client: client,
		workqueue: workqueue,
	}
}


func (c *Controller) Run() {
	wait.Until(c.RunOnce, time.Second, wait.NeverStop)
}


func (c *Controller) RunOnce() {
	tasks := make([]*dosv1.Task, 0)
	timer := time.NewTimer(1 * time.Second)
	run := func() {
		select {
		case <-timer.C:
			return
		default:
			for {
				if c.workqueue.Len() == 0 {
					return
				}
				taskObj, shutdown := c.workqueue.Get()
				if shutdown {
					log.Errorf("queue shutdown")
				}
				task := &dosv1.Task{}
				taskByte := taskObj.(string)
				err := json.Unmarshal([]byte(taskByte), &task)
				if err != nil {
					log.Errorf("error unmarshal to task")
				}
				tasks = append(tasks, task)
			}
		}
	}
	run()
	if len(tasks) == 0 {
		return
	}
	//ss := snapshot.NewSnapshot(c.client)
	//// snapshot := snapshot.MockSnapshot(c.client)
	//alg := algorithm.GetBatchScheduleAlgorithm("batch_max_throughput")
	//result := alg.Schedule(tasks, ss)
	//fmt.Println(result)
	// 加env
	result := MockScheduleResult()
	for _, task := range tasks {
		taskName := task.Namespace + "/" + task.Name
		if task.Annotations == nil {
			task.Annotations = make(map[string]string)
		}
		task.Annotations[ScheduleTimeAnnotation] = strconv.Itoa(int(time.Now().Unix()))
		task.Annotations[ScheduleNodeAnnotation] = result[taskName].NodeName
		task.Annotations[ScheduleGPUIDAnnotation] = result[taskName].GpuId[0]
		taskByte, _ := json.Marshal(task)
		_, err := c.client.UpdateResource(string(taskByte))
		if err != nil {
			fmt.Println(err)
		}
	}
	// run pod

}
func (c *Controller) ProcessNextItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
	}
	fmt.Println("get a obj from the queue ", obj)
	return true
}


func MockScheduleResult()map[string]algorithm.ScheduleResult {
	result := make(map[string]algorithm.ScheduleResult)
	result["default/task1"] = algorithm.ScheduleResult {
		GpuId:    []string{"GPU-" + uuid.New().String()},
		NodeName: "dell04",
	}
	result["default/task2"] = algorithm.ScheduleResult {
		GpuId:    []string{"GPU-" + uuid.New().String()},
		NodeName: "dell04",
	}
	result["default/task3"] = algorithm.ScheduleResult {
		GpuId:    []string{"GPU-" + uuid.New().String()},
		NodeName: "dell04",
	}
	return result
}