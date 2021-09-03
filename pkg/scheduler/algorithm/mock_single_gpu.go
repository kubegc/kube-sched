/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	"github.com/kubesys/kubernetes-scheduler/pkg/scheduler"
)

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/

type MockSingleGPU struct {
	scheduler.Algorithm
}

func NewMockSingleGPU() *MockSingleGPU {
	return &MockSingleGPU{}
}

func (alg *MockSingleGPU) Schedule(requestMemory, requestCore int, availableNode []string, resourceOnNode map[string]*scheduler.NodeResource) *scheduler.Result {
	result := scheduler.Result{
		NodeName: "dell04",
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
	}
	return &result
}

func (alg *MockSingleGPU) Name() string {
	return "MockSingleGPU"
}
