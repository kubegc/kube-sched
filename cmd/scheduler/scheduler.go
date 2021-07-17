/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package main

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/
import (
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	"github.com/kubesys/kubernetes-scheduler/pkg/scheduler"
	alg "github.com/kubesys/kubernetes-scheduler/pkg/scheduler/algorithm"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
)

var (
	masterURL = "https://124.70.64.232:6443"
	token = ""
)

func main() {
	client := kubesys.NewKubernetesClient(masterURL, token)
	client.Init()

	podMgr := scheduler.NewTaskManager(util.NewLinkedQueue())
	nodeMgr := scheduler.NewNodeManager(util.NewLinkedQueue())
	algorithm := alg.NewMockSingleGPU()

	decider := scheduler.NewDecider(client, podMgr, nodeMgr, algorithm)
	go decider.Run()

	decider.Listen(podMgr, nodeMgr)

}