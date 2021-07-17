/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/
package scheduler

import (
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
)

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/
type Result struct {
	GpuCount map[string]int
	GpuId []string
	NodeName string
}
type Algorithm interface {
	Name() string
	Schedule(task []*jsonutil.ObjectNode, nodes *util.LinkedQueue) map[string]Result
}

