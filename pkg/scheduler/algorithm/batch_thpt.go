package algorithm

import (
	"fmt"
	doslabv1 "github.com/kubesys/kubernetes-scheduler/pkg/apis/doslab.io/v1"
	"github.com/kubesys/kubernetes-scheduler/pkg/scheduler/snapshot"
	"math"
	"sort"
)

type BatchThroughputScheduleAlgorithm struct {}


func (bt *BatchThroughputScheduleAlgorithm) Name() string {
	return "batch_max_throughput"
}

func (bt *BatchThroughputScheduleAlgorithm) Schedule(tasks []*doslabv1.Task, snapshot *snapshot.Snapshot) map[string]ScheduleResult {
	numTask := len(tasks)
	if numTask == 0 {
		return nil
	}
	gpuModels := make(map[string]bool)
	gpuNum := make(map[string]int)
	for _, gpu := range snapshot.GPUs {
		gpuNum[gpu.Spec.Model]++
		gpuModels[gpu.Spec.Model] = true
	}
	// task -> speedups

	speedups := make(map[string]map[string]float64)

	fairNum := make(map[string]map[string]int)
	actualNum := make(map[string]map[string]int)

	fairThroughput := make(map[string]float64)
	actualThroughput := make(map[string]float64)

	for _, task := range tasks {
		taskName := GetName(task)
		if fairNum[taskName] == nil {
			fairNum[taskName] = make(map[string]int)
		}
		for model, num := range gpuNum {
			fairNum[taskName][model] = num / numTask
		}
		if speedups[taskName] == nil {
			speedups[taskName] = make(map[string]float64)
		}
		speedups[taskName]["Tesla K80"] = 1.0
		speedups[taskName]["Tesla V100"] = GetSpeedup(task)
		fairThroughput[taskName] = 0.0
		for model, num := range fairNum[taskName] {
			fairThroughput[taskName] += float64(num) * speedups[taskName][model]
		}
	}
	sort.Slice(tasks, func(i, j int) bool {
		speedupi := tasks[i].Annotations["speedup"]
		speedupj := tasks[j].Annotations["speedup"]
		return speedupi < speedupj
	})
	victims := tasks[:len(tasks) - 1]
	victor := tasks[len(tasks) - 1]

	for _, task := range victims {
		taskName := GetName(task)
		models := make([]string, 0)
		for gpuModel, _ := range gpuModels {
			models = append(models, gpuModel)
		}
		sort.Slice(models, func(i, j int) bool {
			return speedups[taskName][models[i]] < speedups[taskName][models[j]]
		})

		thpt := fairThroughput[taskName]

		for _, model := range models {
			need := int(math.Ceil(thpt / speedups[taskName][model]))
			if need <= gpuNum[model] {
				gpuNum[model] -= need
				if actualNum[taskName] == nil {
					actualNum[taskName] = make(map[string]int)
				}
				actualNum[taskName][model] = need
				actualThroughput[taskName] += float64(need) * speedups[taskName][model]
				break
			} else {
				if actualNum[taskName] == nil {
					actualNum[taskName] = make(map[string]int)
				}
				actualNum[taskName][model] = gpuNum[model]
				actualThroughput[taskName] += float64(gpuNum[model]) * speedups[taskName][model]
				thpt -= float64(gpuNum[model]) * speedups[taskName][model]
				gpuNum[model] = 0
			}
		}
	}

	victorName := GetName(victor)
	for model, _ := range gpuModels {
		if actualNum[victorName] == nil {
			actualNum[victorName] = make(map[string]int)
		}
		actualNum[victorName][model] = gpuNum[model]
		actualThroughput[victorName] += float64(gpuNum[model]) * speedups[victorName][model]
		gpuNum[model] = 0
	}
	fmt.Println("------Schedule Result------")
	fmt.Println(actualNum)
	fmt.Println("------Throughput------")
	fmt.Println(actualThroughput)
	fmt.Println(gpuNum)
	return nil
}
