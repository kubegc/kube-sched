package algorithm

import (
	doslabv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler/snapshot"
)

type ScheduleResult struct {
	GpuCount map[string]int
	GpuId []string
	NodeName string
}
type ScheduleAlgorithm interface {
	Name() string
}


type SingleScheduleAlgorithm interface {
	ScheduleAlgorithm
	Schedule(task *doslabv1.Task, snapshot *snapshot.Snapshot) ScheduleResult
}

type BatchScheduleAlgorithm interface {
	ScheduleAlgorithm
	Schedule(tasks []*doslabv1.Task, snapshot *snapshot.Snapshot) map[string]ScheduleResult
}



