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
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
)

var (
	masterURL = "https://124.70.64.232:6443"
	token = ""
)

func main() {
	client := kubesys.NewKubernetesClient(masterURL, token)
	client.Init()

	queue := util.NewLinkedQueue()

	decider := scheduler.NewDecider(client, queue)
	go decider.Run()

	decider.Listen(scheduler.NewTaskManager(queue))
}