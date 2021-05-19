package handler
//
//import (
//	"encoding/json"
//	"k8s.io/client-go/util/workqueue"
//)
//
//type TaskHandler struct {
//	workqueue workqueue.RateLimitingInterface
//	gpuPod2Request map[string]map[string]string
//	gpuPod2Limit map[string]map[string]string
//	gpuPod2Memory map[string]map[string]string
//	gpuPod2Port map[string]map[string]string
//}
//
//
//func NewTaskHandler(wq workqueue.RateLimitingInterface) *TaskHandler {
//	return &TaskHandler{
//		workqueue: wq,
//		gpuPod2Request: make(map[string]map[string]string),
//		gpuPod2Limit: make(map[string]map[string]string),
//		gpuPod2Memory: make(map[string]map[string]string),
//		gpuPod2Port: make(map[string]map[string]string),
//	}
//}
//
//
//func (th *TaskHandler) DoAdded(obj map[string]interface{}) {
//	jb, _ := json.Marshal(obj)
//	th.workqueue.Add(string(jb))
//}
//
//func (th *TaskHandler) DoModified(obj map[string]interface{}) {
//	jb, _ := json.Marshal(obj)
//	th.workqueue.Add(string(jb))
//}
//
//func (th *TaskHandler) DoDeleted(obj map[string]interface{}) {
//	jb, _ := json.Marshal(obj)
//	th.workqueue.Add(string(jb))
//}