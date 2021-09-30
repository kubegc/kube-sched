/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package main

import (
	"flag"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	"github.com/kubesys/kubernetes-scheduler/pkg/scheduler"
	alg "github.com/kubesys/kubernetes-scheduler/pkg/scheduler/algorithm"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
	log "github.com/sirupsen/logrus"
)

var (
	masterUrl = flag.String("masterUrl", "", "Kubernetes master url.")
	token     = flag.String("token", "", "Kubernetes client token.")
)

func main() {
	flag.Parse()
	client := kubesys.NewKubernetesClient(*masterUrl, *token)
	client.Init()

	log.Infoln("Starting pod scheduler.")

	podMgr := scheduler.NewPodManager(util.NewLinkedQueue(), util.NewLinkedQueue())
	gpuMgr := scheduler.NewGpuManager(util.NewLinkedQueue())
	nodeMgr := scheduler.NewNodeManager(util.NewLinkedQueue())
	algorithm := alg.NewMockSingleGPU()

	decider := scheduler.NewDecider(client, podMgr, gpuMgr, nodeMgr, algorithm)
	decider.Listen(podMgr, gpuMgr, nodeMgr)

	decider.Run()
}
