package algorithm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	doslabv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler/snapshot"
	"sort"
	"strconv"
)

type BatchTradeScheduleAlgorithm struct {
}


func (fb *BatchTradeScheduleAlgorithm) Name() string {
	return "batch_trade"
}

func (fb *BatchTradeScheduleAlgorithm) Schedule(tasks []*doslabv1.Task, snapshot *snapshot.Snapshot) map[string]ScheduleResult {
	numTask := len(tasks)
	if numTask == 0 {
		return nil
	}

	sort.Slice(tasks, func(i, j int) bool {
		speedupi := tasks[i].Annotations["speedup"]
		speedupj := tasks[j].Annotations["speedup"]
		return speedupi > speedupj
	})
	gpuNum := make(map[string]int)
	for _, gpu := range snapshot.GPUs {
		gpuNum[gpu.Spec.Model]++
	}
	// task -> speedups
	speedups := make(map[string]map[string]float64)

	fairNum := make(map[string]map[string]int)
	tradeNum := make(map[string]map[string]int)

	fairThroughput := make(map[string]float64)
	tradeThroughput := make(map[string]float64)
	for _, task := range tasks {
		taskName := GetName(task)
		if fairNum[taskName] == nil {
			fairNum[taskName] = make(map[string]int)
		}

		if tradeNum[taskName] == nil {
			tradeNum[taskName] = make(map[string]int)
		}
		for model, num := range gpuNum {
			fairNum[taskName][model] = num / numTask
			tradeNum[taskName][model] = num / numTask
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
	fmt.Println(fairThroughput)
	trader := GetName(tasks[0])
	tradee := GetName(tasks[len(tasks) - 1])
	price := int(GetSpeedup(tasks[len(tasks) - 2]))

	tradeCount := Min(fairNum[trader]["Tesla K80"] / price, fairNum[tradee]["Tesla V100"])
	fmt.Println(tradeCount)
	// trader
	tradeNum[trader]["Tesla K80"] -= tradeCount * price
	tradeNum[trader]["Tesla V100"] += tradeCount
	//tradee
	tradeNum[tradee]["Tesla K80"] += tradeCount * price
	tradeNum[tradee]["Tesla V100"] -= tradeCount

	fmt.Println(tradeNum)
	
	for _, task := range tasks {
		taskName := GetName(task)
		tradeThroughput[taskName] = 0.0
		for model, num := range tradeNum[taskName] {
			tradeThroughput[taskName] += float64(num) * speedups[taskName][model]
		}
	}
	fmt.Println(tradeThroughput)
	return nil

}


func GetName(task *doslabv1.Task) string {
	return task.Namespace + "/" + task.Name
}

func GetSpeedup(task *doslabv1.Task) float64 {
	speedup, err := strconv.ParseFloat(task.Annotations["speedup"], 64)
	if err != nil {
		log.Errorf("error parse task speedup, %s", err)
		return 0.0
	} else {
		return speedup
	}
}

func Min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}