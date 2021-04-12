package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type Snapshot struct {
	client *kubesys.KubernetesClient
	Nodes map[string]*v1.Node
	GPUs map[string]*GPU
}

func NewSnapshot(client *kubesys.KubernetesClient) *Snapshot {
	cache := &Snapshot{
		client: client,
		Nodes:  make(map[string]*v1.Node),
	}
	objNode, err := client.ListResources("Node", "")
	if err != nil {
		log.Errorf("error list nodes, %s", err)
	}
	jb, err := json.Marshal(objNode.Object)
	if err != nil {
		log.Errorf("error marshal node, %s", err)
	}
	var nodelist v1.NodeList
	err = json.Unmarshal(jb, &nodelist)
	fmt.Println(nodelist.Items[0].GetName())
	return cache
}


