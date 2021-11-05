/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

import "github.com/kubesys/kubernetes-scheduler/pkg/scheduler"

type ForFlyingFish struct {
	scheduler.Algorithm
}

func NewForFlyingFish() *ForFlyingFish {
	return &ForFlyingFish{}
}

func (alg *ForFlyingFish) Schedule(requestMemory, requestCore int64, availableNode []string, resourceOnNode map[string]*scheduler.NodeResource) *scheduler.Result {
	result := scheduler.Result{
		NodeName: "dell04",
		GpuUuid:  []string{"GPU-21f591ed-d77b-3a27-c674-51375d2e4fd9"},
	}
	return &result
}

func (alg *ForFlyingFish) Name() string {
	return "MockSingleGPU"
}