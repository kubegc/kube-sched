/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	"encoding/json"
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
	log "github.com/sirupsen/logrus"
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
func (taskMgr *TaskManager) DoAdded(obj map[string]interface{}) {
	taskMgr.queue.Add(jsonutil.NewObjectNodeWithValue(obj))
	jb, _ := json.Marshal(obj)
	log.Info("adding task: " + string(jb))
}

//TODO
func (taskMgr *TaskManager) DoModified(obj map[string]interface{}) {
	//
}

//TODO
func (taskMgr *TaskManager) DoDeleted(obj map[string]interface{}) {

}