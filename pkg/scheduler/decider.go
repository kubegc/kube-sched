/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	"encoding/json"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
	log "github.com/sirupsen/logrus"
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

func NewDecider(client *kubesys.KubernetesClient, queue *util.LinkedQueue) *Decider {
	return &Decider{
		Client: client,
		Queue:  queue,
	}
}

func (d *Decider) Run() {
	for {
		if d.Queue.Len() == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		d.runOnce()
	}
}

func (d *Decider) runOnce() {

	tasks := make([]*jsonutil.ObjectNode, 0)

	for {
		if d.Queue.Len() == 0 {
			break
		}

		tasks = append(tasks, d.Queue.Remove())
	}

	if len(tasks) == 0 {
		return
	}

	// 加env
	result := MockScheduleResult()
	for _, task := range tasks {

		meta := task.GetObjectNode("metadata")

		taskName := meta.GetString("namespace")
		taskName += "/" + meta.GetString("name")

		labels := jsonutil.NewObjectNodeWithValue(make(map[string]interface{}))

		if meta.Object["labels"] != nil {
			labels = meta.GetObjectNode("labels")
		}

		labels.Object[ScheduleNodeAnnotation] = result[taskName].NodeName
		labels.Object[ScheduleGPUIDAnnotation] = result[taskName].GpuId[0]

		meta.Object["labels"] = labels.Object

		taskByte, _ := json.Marshal(task.Object)

		_, err := d.Client.UpdateResource(string(taskByte))

		if err != nil {
			log.Error("Fail to schedule " + taskName)
			log.Error(err)
		} else {
			log.Info("Scheduling " + taskName + " successful")
		}
	}

}


func MockScheduleResult() map[string] ScheduleResult {
	result := make(map[string]ScheduleResult)
	result["default/task1"] = ScheduleResult {
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	result["default/task2"] = ScheduleResult {
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	result["default/task3"] = ScheduleResult {
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	return result
}

func (d *Decider) Listen(taskMgr *TaskManager) {
	watcher := kubesys.NewKubernetesWatcher(d.Client, taskMgr)
	d.Client.WatchResources("Task", "default", watcher)
}