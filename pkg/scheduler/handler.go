package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/kubesys/kubernetes-scheduler/pkg/util"
)

type WorkQueueHandler struct {
	workqueue *util.LinkedQueue
}

func NewWorkQueueHandler(workqueue *util.LinkedQueue) *WorkQueueHandler {
	return &WorkQueueHandler{workqueue: workqueue}
}

func (w *WorkQueueHandler) DoAdded(obj map[string]interface{}) {
	jb, _ := json.Marshal(obj)
	// fmt.Println(string(jb))
	fmt.Println("adding task")
	w.workqueue.Add(string(jb))
}

func (w *WorkQueueHandler) DoModified(obj map[string]interface{}) {
	fmt.Println(obj)
}

func (w *WorkQueueHandler) DoDeleted(obj map[string]interface{}) {
	fmt.Println(obj)
}