package handler

import (
	"encoding/json"
	"k8s.io/client-go/util/workqueue"
)

type TaskHandler struct {
	workqueue workqueue.RateLimitingInterface
}


func NewTaskHandler(wq workqueue.RateLimitingInterface) *TaskHandler {
	return &TaskHandler{
		workqueue: wq,
	}
}


func (th *TaskHandler) DoAdded(obj map[string]interface{}) {
	jb, _ := json.Marshal(obj)
	th.workqueue.Add(string(jb))
}

func (th *TaskHandler) DoModified(obj map[string]interface{}) {
	jb, _ := json.Marshal(obj)
	th.workqueue.Add(string(jb))
}

func (th *TaskHandler) DoDeleted(obj map[string]interface{}) {
	jb, _ := json.Marshal(obj)
	th.workqueue.Add(string(jb))
}