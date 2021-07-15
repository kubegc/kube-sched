package scheduler

import (
	v1 "github.com/kubesys/kubernetes-scheduler/pkg/apis/doslab.io/v1"
	"sort"
)

func Scheduled(task *v1.Task) bool {

	if task.Annotations == nil {
		return false
	}

	return task.Annotations[ScheduleTimeAnnotation] != "" &&
		task.Annotations[ScheduleNodeAnnotation] != "" &&
		task.Annotations[ScheduleGPUIDAnnotation] != ""
}

func Compare(v1 []string, v2 []string) bool {
	if (v1 == nil && v2 == nil) || (v1 != nil && v2 == nil) || (v1 == nil && v2 != nil) || len(v1) != len(v2) {
		return false
	}
	sort.Strings(v1)
	sort.Strings(v2)
	for i := 0; i < len(v1); i++ {
		if v1[i] != v2[i] {
			return false
		}
	}
	return true
}

func SortedKeys(m map[string]bool) []string{
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}




