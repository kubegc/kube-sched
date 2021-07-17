/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	"encoding/json"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"time"
)

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/

type Decider struct {
	Client  *kubesys.KubernetesClient
	PodMgr  *TaskManager
	NodeMgr *NodeManager
	Algorithm interface{}
}

func NewDecider(client *kubesys.KubernetesClient, podMgr *TaskManager, nodeMgr *NodeManager, algorithm interface{}) *Decider {
	return &Decider{
		Client:  client,
		PodMgr:  podMgr,
		NodeMgr: nodeMgr,
		Algorithm: algorithm,
	}
}

func (decider *Decider) Run() {
	for {
		if decider.PodMgr.queue.Len() == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		decider.runOnce()
	}
}

func (decider *Decider) runOnce() {

	tasks := make([]*jsonutil.ObjectNode, 0)

	for {
		if decider.PodMgr.queue.Len() == 0 {
			break
		}

		task := decider.PodMgr.queue.Remove()

		if task != nil {
			tasks = append(tasks, task)
		}
	}

	if len(tasks) == 0 {
		return
	}

	// 加env
	result := decider.Algorithm.(Algorithm).Schedule(tasks, decider.NodeMgr.queue)
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

		_, err := decider.Client.UpdateResource(string(taskByte))

		if err != nil {
			log.Error("Fail to schedule " + taskName)
			log.Error(err)
		} else {
			log.Info("Scheduling " + taskName + " successful")
		}
	}

}


func (decider *Decider) Listen(taskMgr *TaskManager, nodeMgr *NodeManager) {

	taskWatcher := kubesys.NewKubernetesWatcher(decider.Client, taskMgr)
	go decider.Client.WatchResources("Task", "", taskWatcher)

	nodeWatcher := kubesys.NewKubernetesWatcher(decider.Client, nodeMgr)
	decider.Client.WatchResources("Node", "", nodeWatcher)
}