/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

import (
	"github.com/kubesys/kubernetes-scheduler/pkg/scheduler"
	"time"
)

type MockSingleGPU struct {
	scheduler.Algorithm
}

func NewMockSingleGPU() *MockSingleGPU {
	return &MockSingleGPU{}
}

func (alg *MockSingleGPU) Schedule(requestMemory, requestCore int, availableNode []string, resourceOnNode map[string]*scheduler.NodeResource) *scheduler.Result {
	gpuUuid := ""
	if time.Now().Nanosecond() % 2 == 1 {
		gpuUuid = "GPU-21f591ed-d77b-3a27-c674-51375d2e4fd9"
	} else {
		gpuUuid = " GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"
	}

	result := scheduler.Result{
		NodeName: "dell04",
		GpuUuid:  []string{gpuUuid},
	}
	return &result
}

func (alg *MockSingleGPU) Name() string {
	return "MockSingleGPU"
}
