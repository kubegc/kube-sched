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

const (
	DefaultMasterUrl = "https://133.133.135.42:6443"
	DefaultToken     = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjN2U3VZUk16R3ZfZGNaMkw4bVktVGlRWnJGZFB2NWprU1lrd0hObnNBVFEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbi10Z202ZyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjI0MjNlMDJmLTdmYzAtNDEzYi04ODczLTc0YTM3MTFkMzdkOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.KVJ7NC4NAWViLy2YkFFzzg0G4NcKnAZzw8VYooyXaLQlyfJWysR0giU8QLcSRs5BqIagff2EcVBuVHmSE4o1Zt3AMayStk-stwdtQre28adKYwR4aJLtfa1Wqmw--RiBHZmOjOmzynDdtWEe_sJPl4bGSxMvjFEKy6OepXOctnqZjUq4x2mMK-FID5hmeoHY6oAcfrRuAJsHRuLEAJQzLiMAf9heTuRNxcv3OTyfGtLOOj9risr59wilC_JWVPC5DC5TkEe4-8OeWg_mKA-lwSss_nyGMCsBqPIdPeyd3RQQ9ADPDq-JP2Nci0zoqOEwgZu3nQ3wOovR7lFBbRxsQQ"
)

var (
	masterUrl = flag.String("masterUrl", DefaultMasterUrl, "Kubernetes master url.")
	token     = flag.String("token", DefaultToken, "Kubernetes client token.")
)

func main() {
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
