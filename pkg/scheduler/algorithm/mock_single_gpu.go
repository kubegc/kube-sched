/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	"github.com/kubesys/kubernetes-scheduler/pkg/scheduler"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
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

func (alg *MockSingleGPU) Schedule(tasks []*jsonutil.ObjectNode, nodes *util.LinkedQueue) map[string]scheduler.Result  {
	result := make(map[string]scheduler.Result)
	result["default/task1"] = scheduler.Result{
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	result["default/task2"] = scheduler.Result{
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	result["default/task3"] = scheduler.Result{
		GpuId:    []string{"GPU-da33250c-6bee-6f8d-dd97-f1aa43d95783"},
		NodeName: "dell04",
	}
	return result
}

func (alg *MockSingleGPU) Name() string  {
	return "MockSingleGPU"
}
