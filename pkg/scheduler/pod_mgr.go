/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

import (
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
	"sync"
)

type PodManager struct {
	queueOfAdded   *util.LinkedQueue
	queueOfDeleted *util.LinkedQueue
	muOfAdd        sync.Mutex
	muOfDelete     sync.Mutex
}

func NewPodManager(queueOfAdded, queueOfDeleted *util.LinkedQueue) *PodManager {
	return &PodManager{queueOfAdded: queueOfAdded, queueOfDeleted: queueOfDeleted}
}

func (podMgr *PodManager) DoAdded(obj map[string]interface{}) {
	podMgr.muOfAdd.Lock()
	podMgr.queueOfAdded.Add(jsonutil.NewObjectNodeWithValue(obj))
	podMgr.muOfAdd.Unlock()
}

func (podMgr *PodManager) DoModified(obj map[string]interface{}) {
}

func (podMgr *PodManager) DoDeleted(obj map[string]interface{}) {
	podMgr.muOfDelete.Lock()
	podMgr.queueOfDeleted.Add(jsonutil.NewObjectNodeWithValue(obj))
	podMgr.muOfDelete.Unlock()
}
