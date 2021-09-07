/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package node_daemon

import (
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	"time"
)

type NodeDaemon struct {
	Client   *kubesys.KubernetesClient
	PodMgr   *PodManager
	NodeName string
}

func NewNodeDaemon(client *kubesys.KubernetesClient, podMgr *PodManager, nodeName string) *NodeDaemon {
	return &NodeDaemon{
		Client:   client,
		PodMgr:   podMgr,
		NodeName: nodeName,
	}
}

func (daemon *NodeDaemon) Run() {

	for {

		if daemon.PodMgr.queueOfAdded.Len() > 0 {
			daemon.PodMgr.muOfAdd.Lock()
			pod := daemon.PodMgr.queueOfAdded.Remove()
			daemon.PodMgr.muOfAdd.Unlock()
			time.Sleep(5 * time.Millisecond)
			go daemon.addPod(pod)
		}

		if daemon.PodMgr.queueOfDeleted.Len() > 0 {
			daemon.PodMgr.muOfDelete.Lock()
			pod := daemon.PodMgr.queueOfDeleted.Remove()
			daemon.PodMgr.muOfDelete.Unlock()
			time.Sleep(5 * time.Millisecond)
			go daemon.deletePod(pod)
		}

	}
}

func (daemon *NodeDaemon) Listen(podMgr *PodManager) {
	podWatcher := kubesys.NewKubernetesWatcher(daemon.Client, podMgr)
	daemon.Client.WatchResources("Pod", "", podWatcher)
}

func (daemon *NodeDaemon) addPod(pod *jsonutil.ObjectNode) {

}

func (daemon *NodeDaemon) deletePod(pod *jsonutil.ObjectNode) {

}
