/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

import (
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
	"sync"
)

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/

type NodeManager struct {
	queue *util.LinkedQueue
	mu    sync.Mutex
}

func NewNodeManager(queue *util.LinkedQueue) *NodeManager {
	return &NodeManager{queue: queue}
}

func (nodeMgr *NodeManager) DoAdded(obj map[string]interface{}) {
}

func (nodeMgr *NodeManager) DoModified(obj map[string]interface{}) {
	nodeMgr.mu.Lock()
	nodeMgr.queue.Add(jsonutil.NewObjectNodeWithValue(obj))
	nodeMgr.mu.Unlock()
}

func (nodeMgr *NodeManager) DoDeleted(obj map[string]interface{}) {
}
