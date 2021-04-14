package algorithm

import (
	doslabv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler/snapshot"
)

type BatchFairScheduleAlgorithm struct {

}


func (fb *BatchFairScheduleAlgorithm) Name() string {
	return "batch_fair"
}

func (fb *BatchFairScheduleAlgorithm) Schedule(tasks []*doslabv1.Task, snapshot *snapshot.Snapshot) map[string]ScheduleResult {
	return nil
}
