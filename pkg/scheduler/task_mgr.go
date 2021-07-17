/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
)

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/

type TaskManager struct {
	queue *util.LinkedQueue
}

func NewTaskManager(queue *util.LinkedQueue) *TaskManager {
	return &TaskManager{queue: queue}
}

//TODO
func (w *TaskManager) DoAdded(obj map[string]interface{}) {
	jb, _ := json.Marshal(obj)
	w.queue.Add(string(jb))
	fmt.Println("adding task: " + string(jb))
}

//TODO
func (w *TaskManager) DoModified(obj map[string]interface{}) {
	fmt.Println(obj)
}

//TODO
func (w *TaskManager) DoDeleted(obj map[string]interface{}) {
	fmt.Println(obj)
}