/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

import (
	"encoding/json"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

/**
 *   authors: yangchen19@otcaix.iscas.ac.cn
 *            wuheng@iscas.ac.cn
 *
 **/

type Decider struct {
	Client         *kubesys.KubernetesClient
	PodMgr         *PodManager
	GpuMgr         *GpuManager
	NodeMgr        *NodeManager
	Algorithm      interface{}
	resourceOnNode map[string]*NodeResource
	gpuUuidToName  map[string]string
	mu             sync.Mutex
}

func NewDecider(client *kubesys.KubernetesClient, podMgr *PodManager, gpuMgr *GpuManager, nodeMgr *NodeManager, algorithm interface{}) *Decider {
	return &Decider{
		Client:         client,
		PodMgr:         podMgr,
		GpuMgr:         gpuMgr,
		NodeMgr:        nodeMgr,
		Algorithm:      algorithm,
		resourceOnNode: make(map[string]*NodeResource),
		gpuUuidToName:  make(map[string]string),
	}
}

func (decider *Decider) Run() {
	for {
		if decider.PodMgr.queueOfAdded.Len() > 0 {
			decider.PodMgr.muOfAdd.Lock()
			pod := decider.PodMgr.queueOfAdded.Remove()
			decider.PodMgr.muOfAdd.Unlock()
			go decider.addPod(pod)
			time.Sleep(5 * time.Millisecond)
		}

		if decider.PodMgr.queueOfDeleted.Len() > 0 {
			decider.PodMgr.muOfDelete.Lock()
			pod := decider.PodMgr.queueOfDeleted.Remove()
			decider.PodMgr.muOfDelete.Unlock()
			go decider.deletePod(pod)
			time.Sleep(5 * time.Millisecond)
		}

		if decider.GpuMgr.queue.Len() > 0 {
			decider.GpuMgr.mu.Lock()
			gpu := decider.GpuMgr.queue.Remove()
			decider.GpuMgr.mu.Unlock()
			go decider.addGpu(gpu)
			time.Sleep(5 * time.Millisecond)
		}

		if decider.NodeMgr.queue.Len() > 0 {
			decider.NodeMgr.mu.Lock()
			node := decider.NodeMgr.queue.Remove()
			decider.NodeMgr.mu.Unlock()
			go decider.modifyNode(node)
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (decider *Decider) Listen(podMgr *PodManager, gpuMgr *GpuManager, nodeMgr *NodeManager) {

	podWatcher := kubesys.NewKubernetesWatcher(decider.Client, podMgr)
	go decider.Client.WatchResources("Pod", "", podWatcher)

	gpuWatcher := kubesys.NewKubernetesWatcher(decider.Client, gpuMgr)
	go decider.Client.WatchResources("GPU", "", gpuWatcher)

	nodeWatcher := kubesys.NewKubernetesWatcher(decider.Client, nodeMgr)
	decider.Client.WatchResources("Node", "", nodeWatcher)
}

func (decider *Decider) addPod(pod *jsonutil.ObjectNode) {
	spec := pod.GetObjectNode("spec")
	schedulerName := spec.GetString("schedulerName")
	if schedulerName != SchedulerName {
		return
	}

	meta := pod.GetObjectNode("metadata")
	podName := meta.GetString("name")
	namespace := meta.GetString("namespace")

	requestMemory, requestCore := 0, 0
	containers := spec.GetArray("containers")
	for _, c := range containers {
		container := c.(map[string]interface{})
		if _, ok := container["resources"]; !ok {
			continue
		}
		resources := container["resources"].(map[string]interface{})
		if _, ok := resources["limits"]; !ok {
			continue
		}
		limits := resources["limits"].(map[string]interface{})
		if val, ok := limits[ResourceMemory]; ok {
			m, _ := strconv.Atoi(val.(string))
			requestMemory += m
		}
		if val, ok := limits[ResourceCore]; ok {
			m, _ := strconv.Atoi(val.(string))
			requestCore += m
		}
	}

	log.Infof("Scheduling node and gpus for pod %s on namespace %s, which need memory %dMiB and core %d%%.", podName, namespace, requestMemory, requestCore)

	decider.mu.Lock()
	defer decider.mu.Unlock()

	var availableNode []string
	for _, v := range decider.resourceOnNode {
		if v.hasDevicePlugin {
			availableNode = append(availableNode, v.nodeName)
		}
	}
	// TODO: Filter more unavailable nodes.

	result := decider.Algorithm.(Algorithm).Schedule(requestMemory, requestCore, availableNode, decider.resourceOnNode)

	if result == nil {
		log.Warningf("There is no suitable resource for pod %s on namespace %s, try again later.", podName, namespace)
		time.Sleep(100 * time.Millisecond)
		decider.PodMgr.muOfAdd.Lock()
		decider.PodMgr.queueOfAdded.Add(pod)
		decider.PodMgr.muOfAdd.Unlock()
		return
	}

	// Add annotations and bind node
	annotations := &jsonutil.ObjectNode{}
	if meta.Object["annotations"] != nil {
		annotations = meta.GetObjectNode("annotations")
	}
	annotations.Object[AnnAssumeTime] = strconv.FormatInt(time.Now().UnixNano(), 10)
	annotations.Object[AnnAssignedFlag] = "false"
	annotations.Object[ResourceUUID] = result.GpuUuid[0]
	meta.Object["annotations"] = annotations.Object

	podByte, _ := json.Marshal(pod.Object)
	_, err := decider.Client.UpdateResource(string(podByte))
	if err != nil {
		log.Warningf("Failed to add annotations for pod %s on namespace %s, %s, try again later.", podName, namespace, err)
		time.Sleep(100 * time.Millisecond)
		decider.PodMgr.muOfAdd.Lock()
		decider.PodMgr.queueOfAdded.Add(pod)
		decider.PodMgr.muOfAdd.Unlock()
		return
	}

	bind := map[string]interface{}{}
	bind["apiVersion"] = "v1"
	bind["kind"] = "Binding"
	bind["metadata"] = map[string]string{
		"name":      podName,
		"namespace": namespace,
	}
	bind["target"] = map[string]string{
		"apiVersion": "v1",
		"kind":       "Node",
		"name":       result.NodeName,
	}

	bindByte, _ := json.Marshal(bind)
	_, err = decider.Client.CreateResource(string(bindByte))
	if err != nil && err.Error() != "request status 201 Created" {
		log.Warningf("Failed to bind node for pod %s on namespace %s, %s, try again later.", podName, namespace, err)
		time.Sleep(100 * time.Millisecond)
		decider.PodMgr.muOfAdd.Lock()
		decider.PodMgr.queueOfAdded.Add(pod)
		decider.PodMgr.muOfAdd.Unlock()
		return
	}

	// Update resource and GPU CRD
	decider.resourceOnNode[result.NodeName].gpusByUuid[result.GpuUuid[0]].memoryAllocated += requestMemory
	decider.resourceOnNode[result.NodeName].gpusByUuid[result.GpuUuid[0]].coreAllocated += requestCore

	gpuName := decider.gpuUuidToName[result.GpuUuid[0]]
	gpu, err := decider.Client.GetResource("GPU", GPUNamespace, gpuName)
	if err != nil {
		log.Fatalf("Failed to get GPU CRD, %s.", err)
	}
	status := gpu.GetObjectNode("status")
	allocated := status.GetObjectNode("allocated")

	oldMemory, _ := strconv.Atoi(allocated.Object["memory"].(string))
	oldCore, _ := strconv.Atoi(allocated.Object["core"].(string))

	allocated.Object["memory"] = strconv.Itoa(oldMemory + requestMemory)
	allocated.Object["core"] = strconv.Itoa(oldCore + requestCore)

	gpuByte, _ := json.Marshal(gpu.Object)
	_, err = decider.Client.UpdateResource(string(gpuByte))
	if err != nil {
		log.Fatalf("Failed to update GPU CRD, %s.", err)
	}

	log.Infof("Pod %s on namespace %s will run on node %s with %d gpu(s).", podName, namespace, result.NodeName, len(result.GpuUuid))

}

func (decider *Decider) deletePod(pod *jsonutil.ObjectNode) {
	spec := pod.GetObjectNode("spec")
	schedulerName := spec.GetString("schedulerName")
	nodeName := spec.GetString("nodeName")
	if schedulerName != SchedulerName {
		return
	}

	meta := pod.GetObjectNode("metadata")
	podName := meta.GetString("name")
	namespace := meta.GetString("namespace")
	annotations := meta.GetObjectNode("annotations")
	gpuUuid := annotations.GetString(ResourceUUID)

	requestMemory, requestCore := 0, 0
	containers := spec.GetArray("containers")
	for _, c := range containers {
		container := c.(map[string]interface{})
		if _, ok := container["resources"]; !ok {
			continue
		}
		resources := container["resources"].(map[string]interface{})
		if _, ok := resources["limits"]; !ok {
			continue
		}
		limits := resources["limits"].(map[string]interface{})
		if val, ok := limits[ResourceMemory]; ok {
			m, _ := strconv.Atoi(val.(string))
			requestMemory += m
		}
		if val, ok := limits[ResourceCore]; ok {
			m, _ := strconv.Atoi(val.(string))
			requestCore += m
		}
	}

	log.Infof("Releasing resources for pod %s on namespace %s, which need memory %dMiB and core %d%%.", podName, namespace, requestMemory, requestCore)

	decider.mu.Lock()
	defer decider.mu.Unlock()

	// Update resource and GPU CRD
	decider.resourceOnNode[nodeName].gpusByUuid[gpuUuid].memoryAllocated -= requestMemory
	decider.resourceOnNode[nodeName].gpusByUuid[gpuUuid].coreAllocated -= requestCore

	gpuName := decider.gpuUuidToName[gpuUuid]
	gpu, err := decider.Client.GetResource("GPU", GPUNamespace, gpuName)
	if err != nil {
		log.Fatalf("Failed to get GPU CRD, %s.", err)
	}
	status := gpu.GetObjectNode("status")
	allocated := status.GetObjectNode("allocated")

	oldMemory, _ := strconv.Atoi(allocated.Object["memory"].(string))
	oldCore, _ := strconv.Atoi(allocated.Object["core"].(string))

	allocated.Object["memory"] = strconv.Itoa(oldMemory - requestMemory)
	allocated.Object["core"] = strconv.Itoa(oldCore - requestCore)

	gpuByte, _ := json.Marshal(gpu.Object)
	_, err = decider.Client.UpdateResource(string(gpuByte))
	if err != nil {
		log.Fatalf("Failed to update GPU CRD, %s.", err)
	}

	log.Infof("Pod %s on namespace %s is deleled on node %s.", podName, namespace, nodeName)
}

func (decider *Decider) addGpu(gpu *jsonutil.ObjectNode) {
	meta := gpu.GetObjectNode("metadata")
	gpuName := meta.GetString("name")

	spec := gpu.GetObjectNode("spec")
	gpuUuid := spec.GetString("uuid")
	nodeName := spec.GetString("node")

	status := gpu.GetObjectNode("status")
	capacity := spec.GetObjectNode("capacity")
	allocated := status.GetObjectNode("allocated")

	coreCapacity, _ := strconv.Atoi(capacity.Object["core"].(string))
	coreAllocated, _ := strconv.Atoi(allocated.Object["core"].(string))
	memoryCapacity, _ := strconv.Atoi(capacity.Object["memory"].(string))
	memoryAllocated, _ := strconv.Atoi(allocated.Object["memory"].(string))

	hasDevicePlugin := false
	node, err := decider.Client.GetResource("Node", "", nodeName)
	if err != nil {
		log.Fatalf("Failed to get GPU's node, %s.", err)
	}

	nodeStatus := node.GetObjectNode("status")
	nodeCapacity := nodeStatus.GetObjectNode("capacity")
	if val, ok := nodeCapacity.Object[ResourceCore]; ok {
		nodeDeviceCapacity := val.(string)
		if nodeDeviceCapacity != "0" {
			hasDevicePlugin = true
		}
	}

	decider.mu.Lock()
	defer decider.mu.Unlock()

	decider.gpuUuidToName[gpuUuid] = gpuName

	gpuResource := &GpuResource{
		gpuName:         gpuName,
		uuid:            gpuUuid,
		node:            nodeName,
		coreCapacity:    coreCapacity,
		coreAllocated:   coreAllocated,
		memoryCapacity:  memoryCapacity,
		memoryAllocated: memoryAllocated,
	}
	if _, ok := decider.resourceOnNode[nodeName]; ok {
		decider.resourceOnNode[nodeName].gpusByUuid[gpuUuid] = gpuResource
	} else {
		decider.resourceOnNode[nodeName] = &NodeResource{
			nodeName:        nodeName,
			hasDevicePlugin: hasDevicePlugin,
			gpusByUuid:      map[string]*GpuResource{gpuUuid: gpuResource},
		}
	}

	log.Infof("GPU CRD %s, uuid %s added.", gpuName, gpuUuid)
}

func (decider *Decider) modifyNode(node *jsonutil.ObjectNode) {
	meta := node.GetObjectNode("metadata")
	nodeName := meta.GetString("name")

	hasDevicePlugin := false
	nodeStatus := node.GetObjectNode("status")
	nodeCapacity := nodeStatus.GetObjectNode("capacity")
	if val, ok := nodeCapacity.Object[ResourceCore]; ok {
		nodeDeviceCapacity := val.(string)
		if nodeDeviceCapacity != "0" {
			hasDevicePlugin = true
		}
	}

	decider.mu.Lock()
	defer decider.mu.Unlock()

	if val, ok := decider.resourceOnNode[nodeName]; !ok || val.hasDevicePlugin == hasDevicePlugin {
		return
	}

	// Update resource
	decider.resourceOnNode[nodeName].hasDevicePlugin = hasDevicePlugin
	if hasDevicePlugin {
		log.Infof("Node %s now runs device plugin.", nodeName)
	} else {
		log.Infof("Node %s now loses device plugin.", nodeName)
	}
}
