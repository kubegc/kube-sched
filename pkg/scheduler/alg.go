package scheduler

type ScheduleResult struct {
	gpuId string
	nodeName string
}
type ScheduleAlgorithm interface {}

type SingleScheduleAlgorithm interface {
	ScheduleAlgorithm
	Schedule(*doslabv1.Task, snapshot *Snapshot) ScheduleResult
}

type BatchScheduleAlgorithm interface {
	ScheduleAlgorithm
	Schedule(tasks []*doslabv1.Task, snapshot *Snapshot) map[*doslabv1.Task]ScheduleResult
}