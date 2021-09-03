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

	ResourceAssumeTime = "doslab.io/gpu-assume-time"
	ResourceUUID       = "doslab.io/gpu-uuid"
	AnnAssignedFlag    = "doslab.io/gpu-assigned"
)
