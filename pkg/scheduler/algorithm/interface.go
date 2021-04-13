package algorithm

import (
	doslabv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler"
)

type ScheduleResult struct {
	gpuId []string
	nodeName string
}
type ScheduleAlgorithm interface {
	Name() string
}

type SingleScheduleAlgorithm interface {
	ScheduleAlgorithm
	Schedule(task *doslabv1.Task, snapshot *scheduler.Snapshot) ScheduleResult
}

type BatchScheduleAlgorithm interface {
	ScheduleAlgorithm
	Schedule(tasks []*doslabv1.Task, snapshot *scheduler.Snapshot) map[*doslabv1.Task]ScheduleResult
}


