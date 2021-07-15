package algorithm

import (
	"fmt"
	doslabv1 "github.com/kubesys/kubernetes-scheduler/pkg/apis/doslab.io/v1"
	"github.com/kubesys/kubernetes-scheduler/pkg/scheduler/snapshot"
)

type BatchFairScheduleAlgorithm struct {

}


func (bf *BatchFairScheduleAlgorithm) Name() string {
	return "batch_fair"
}

func (bf *BatchFairScheduleAlgorithm) Schedule(tasks []*doslabv1.Task, snapshot *snapshot.Snapshot) map[string]ScheduleResult {
	// 获取各类gpu的数量
	// Model to number
	numTask := len(tasks)
	if numTask == 0 {
		return nil
	}
	gpuNum := make(map[string]int)
	for _, gpu := range snapshot.GPUs {
		gpuNum[gpu.Spec.Model]++
	}
	fairNum := make(map[string]int)

	for model, num := range gpuNum {
		fairNum[model] = num / numTask
	}
	fmt.Println(fairNum)

	scheduleResult := make(map[string]ScheduleResult)

	for _, task := range tasks {
		taskGPUResult := scheduleResult[task.Namespace + "/" + task.Name]
		if taskGPUResult.GpuCount == nil {
			taskGPUResult.GpuCount = make(map[string]int)
		}
		for model, num := range fairNum {
			taskGPUResult.GpuCount[model] = num
		}
	}
	return scheduleResult
}
