package scheduler

import (
	"fmt"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type Controller struct {
	client *kubesys.KubernetesClient
	workqueue workqueue.RateLimitingInterface
}

func NewController(client *kubesys.KubernetesClient, workqueue workqueue.RateLimitingInterface) *Controller {
	return &Controller{
		client: client,
		workqueue: workqueue,
	}
}


func (c *Controller) Run() {
	wait.Until(c.RunOnce, time.Second, wait.NeverStop)
}


func (c *Controller) RunOnce() {
	//timer := time.NewTimer(3 * time.Second)
	//tasks := make(map[string]interface{})

	_ = NewSnapshot(c.client)



	// 加env

	// run pod





}
func (c *Controller) ProcessNextItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
	}
	fmt.Println("get a obj from the queue ", obj)
	return true
}