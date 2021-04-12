package scheduler

import (
	"encoding/json"
	"fmt"
	"k8s.io/client-go/util/workqueue"
)

type WorkQueueHandler struct {
	workqueue workqueue.RateLimitingInterface
}

func NewWorkQueueHandler(workqueue workqueue.RateLimitingInterface) *WorkQueueHandler {
	return &WorkQueueHandler{workqueue: workqueue}
}

func (w *WorkQueueHandler) DoAdded(obj map[string]interface{}) {
	jb, _ := json.Marshal(obj)
	w.workqueue.Add(string(jb))
}

func (w *WorkQueueHandler) DoModified(obj map[string]interface{}) {
	fmt.Println(obj)
}

func (w *WorkQueueHandler) DoDeleted(obj map[string]interface{}) {
	fmt.Println(obj)
}