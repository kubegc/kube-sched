/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

import (
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
	"sync"
)

type GpuResource struct {
	gpuName         string
	uuid            string
	node            string
	coreCapacity    int
	coreAllocated   int
	memoryCapacity  int
	memoryAllocated int
}

type NodeResource struct {
	nodeName        string
	hasDevicePlugin bool
	gpusByUuid      map[string]*GpuResource
}

type GpuManager struct {
	queue *util.LinkedQueue
	mu    sync.Mutex
}

func NewGpuManager(queue *util.LinkedQueue) *GpuManager {
	return &GpuManager{queue: queue}
}

func (gpuMgr *GpuManager) DoAdded(obj map[string]interface{}) {
	gpuMgr.mu.Lock()
	gpuMgr.queue.Add(jsonutil.NewObjectNodeWithValue(obj))
	gpuMgr.mu.Unlock()
}

func (gpuMgr *GpuManager) DoModified(obj map[string]interface{}) {
}

func (gpuMgr *GpuManager) DoDeleted(obj map[string]interface{}) {
}
