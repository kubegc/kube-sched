/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	dosv1 "github.com/kubesys/kubernetes-scheduler/pkg/apis/doslab.io/v1"
	"github.com/kubesys/kubernetes-scheduler/pkg/scheduler/algorithm"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"strconv"
	"time"
)

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/

type Decider struct {
	Client *kubesys.KubernetesClient
	Queue  *util.LinkedQueue
}

func NewDecider(client *kubesys.KubernetesClient) *Decider {
	return &Decider{
		Client: client,
		Queue:  util.NewLinkedQueue(),
	}
}

func (c *Decider) Run() {
	wait.Until(c.RunOnce, time.Second, wait.NeverStop)
}

func (c *Decider) RunOnce() {
	tasks := make([]*dosv1.Task, 0)
	timer := time.NewTimer(1 * time.Second)
	run := func() {
		select {
		case <-timer.C:
			return
		default:
			for {
				if c.Queue.Len() == 0 {
					return
				}
				taskObj  := c.Queue.Remove()
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
		_, err := c.Client.UpdateResource(string(taskByte))
		if err != nil {
			fmt.Println(err)
		}
	}

}
func (c *Decider) ProcessNextItem() bool {
	obj := c.Queue.Remove()
	fmt.Println("get a obj from the Queue ", obj)
	return true
}


func MockScheduleResult()map[string]algorithm.ScheduleResult {
	result := make(map[string]algorithm.ScheduleResult)
	result["default/task1"] = algorithm.ScheduleResult {
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	result["default/task2"] = algorithm.ScheduleResult {
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	result["default/task3"] = algorithm.ScheduleResult {
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	return result
}

func (c *Decider) ListenTask (taskMgr *TaskManager) {
	watcher := kubesys.NewKubernetesWatcher(c.Client, taskMgr)
	c.Client.WatchResources("Task", "default", watcher)
}