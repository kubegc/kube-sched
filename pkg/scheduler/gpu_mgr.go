/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

import (
	"encoding/json"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
	"sync"
)

type GpuResource struct {
	gpuName         string
	uuid            string
	node            string
	coreCapacity    int64
	coreAllocated   int64
	memoryCapacity  int64
	memoryAllocated int64
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
	bytes, _ := json.Marshal(obj)
	gpuMgr.mu.Lock()
	gpuMgr.queue.Add(kubesys.ToJsonObject(bytes))
	gpuMgr.mu.Unlock()
}

func (gpuMgr *GpuManager) DoModified(obj map[string]interface{}) {
}

func (gpuMgr *GpuManager) DoDeleted(obj map[string]interface{}) {
}
