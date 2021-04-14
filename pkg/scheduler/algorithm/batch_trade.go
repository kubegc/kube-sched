package algorithm

import (
	doslabv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler/snapshot"
)

type BatchTradeScheduleAlgorithm struct {

}


func (fb *BatchTradeScheduleAlgorithm) Name() string {
	return "batch_trade"
}

func (fb *BatchTradeScheduleAlgorithm) Schedule(tasks []*doslabv1.Task, snapshot *snapshot.Snapshot) map[string]ScheduleResult {
	return nil
}
