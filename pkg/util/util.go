package util

import (
	v1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
	"kubesys.io/dl-scheduler/pkg/scheduler"
)

func Scheduled(task *v1.Task) bool {

	if task.Annotations == nil {
		return false
	}

	return task.Annotations[scheduler.ScheduleTimeAnnotation] != "" &&
		task.Annotations[scheduler.ScheduleNodeAnnotation] != "" &&
		task.Annotations[scheduler.ScheduleGPUIDAnnotation] != ""
}