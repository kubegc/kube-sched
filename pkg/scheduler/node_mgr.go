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

type NodeManager struct {
	queue *util.LinkedQueue
}

func NewNodeManager(queue *util.LinkedQueue) *NodeManager {
	return &NodeManager{queue: queue}
}

//TODO
func (nodeMgr *NodeManager) DoAdded(obj map[string]interface{}) {
	nodeMgr.queue.Add(jsonutil.NewObjectNodeWithValue(obj))
	jb, _ := json.Marshal(obj)
	log.Info("adding node: " + string(jb))
}

//TODO
func (nodeMgr *NodeManager) DoModified(obj map[string]interface{}) {
	//
}

//TODO
func (nodeMgr *NodeManager) DoDeleted(obj map[string]interface{}) {

}
