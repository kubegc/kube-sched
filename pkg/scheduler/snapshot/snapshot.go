package snapshot

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	dosv1 "kubesys.io/dl-scheduler/pkg/apis/doslab.io/v1"
)

type Snapshot struct {
	client *kubesys.KubernetesClient
	Nodes map[string]*v1.Node
	GPUs map[string]*dosv1.GPU
}

func NewSnapshot(client *kubesys.KubernetesClient) *Snapshot {
	snapshot := &Snapshot{
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
	// fmt.Println(nodelist.Items[0].GetName())
	return snapshot
}

func MockSnapshot(client *kubesys.KubernetesClient) *Snapshot {
	snapshot := &Snapshot{
		client: client,
		Nodes:  make(map[string]*v1.Node),
		GPUs:   make(map[string]*dosv1.GPU),
	}
	for i := 0; i < 60; i++ {
		gpu := &dosv1.GPU{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-gpu-%d", "dell04", i),
			},
			Spec: dosv1.GPUSpec{
				UUID:   "GPU-" + uuid.New().String(),
				Model:  "Tesla K80",
				Family: "Turing",
				Capacity: dosv1.R{
					Core:   "1.0",
					Memory: 10240,
				},
				Node: "dell04",
			},
			Status: dosv1.GPUStatus{
				Allocated: dosv1.R{
					Core:   "0.0",
					Memory: 0,
				},
				PodMap: nil,
			},
		}

		snapshot.GPUs[gpu.Name] = gpu
	}

	for i := 60; i < 72; i++ {
		gpu := &dosv1.GPU{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-gpu-%d", "dell04", i),
			},
			Spec: dosv1.GPUSpec{
				UUID:   "GPU-" + uuid.New().String(),
				Model:  "Tesla V100",
				Family: "Turing",
				Capacity: dosv1.R{
					Core:   "1.0",
					Memory: 10240,
				},
				Node: "dell04",
			},
			Status: dosv1.GPUStatus{
				Allocated: dosv1.R{
					Core:   "0.0",
					Memory: 0,
				},
				PodMap: nil,
			},
		}
		snapshot.GPUs[gpu.Name] = gpu
	}


	return snapshot
}