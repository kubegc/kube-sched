package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	dosv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler/algorithm"
	"kubesys.io/dl-scheduler/pkg/scheduler/snapshot"
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
	_ = snapshot.NewSnapshot(c.client)
	snapshot := snapshot.MockSnapshot(c.client)
	alg := algorithm.GetBatchScheduleAlgorithm("batch_trade")
	result := alg.Schedule(tasks, snapshot)
	fmt.Println(result)
	// 加env

	// run pod





}
func (c *Controller) ProcessNextItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
	}
	fmt.Println("get a obj from the queue ", obj)
	return true
}