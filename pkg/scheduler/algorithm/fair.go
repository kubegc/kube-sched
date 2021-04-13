package algorithm

import (
	doslabv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler"
)

type FairBatchScheduleAlgorithm struct {

}


func (fb *FairBatchScheduleAlgorithm) Name() string {
	return "fair_batch"
}

func (fb *FairBatchScheduleAlgorithm) Schedule(tasks []*doslabv1.Task, snapshot *scheduler.Snapshot) map[*doslabv1.Task]ScheduleResult {
	return nil
}
