/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
)

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/
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
	Schedule(task *[]jsonutil.ObjectNode, nodes *[]jsonutil.ObjectNode) ScheduleResult
}
