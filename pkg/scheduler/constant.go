/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/
const (
	SchedulerName = "doslab-gpu-scheduler"

	GPUNamespace = "default"

	ResourceMemory = "doslab.io/gpu-memory"
	ResourceCore   = "doslab.io/gpu-core"
	ResourceUUID   = "doslab.io/gpu-uuid"

	AnnAssumeTime   = "doslab.io/gpu-assume-time"
	AnnAssignedFlag = "doslab.io/gpu-assigned"
)
