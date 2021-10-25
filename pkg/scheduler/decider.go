/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package scheduler

import (
	"encoding/json"
	jsonObj "github.com/kubesys/kubernetes-client-go/pkg/json"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

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
	go decider.Client.WatchResources("Node", "", nodeWatcher)
}

func (decider *Decider) addPod(pod *jsonObj.JsonObject) {
	spec := pod.GetJsonObject("spec")
	schedulerName, err := spec.GetString("schedulerName")
	if err != nil || schedulerName != SchedulerName {
		return
	}

	meta := pod.GetJsonObject("metadata")
	podName, err := meta.GetString("name")
	if err != nil {
		log.Fatalln("Failed to get pod name.")
	}
	namespace, err := meta.GetString("namespace")
	if err != nil {
		log.Fatalln("Failed to get pod namespace.")
	}

	requestMemory, requestCore := int64(0), int64(0)
	containers := spec.GetJsonArray("containers")
	for _, c := range containers.Values() {
		container := c.JsonObject()
		if !container.HasKey("resources") {
			continue
		}
		resources := container.GetJsonObject("resources")
		if !resources.HasKey("limits") {
			continue
		}
		limits := resources.GetJsonObject("limits")
		if val, err := limits.GetString(ResourceMemory); err == nil {
			m, _ := strconv.ParseInt(val, 10, 64)
			requestMemory += m
		}
		if val, err := limits.GetString(ResourceCore); err == nil {
			m, _ := strconv.ParseInt(val, 10, 64)
			requestCore += int64(m)
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
	annotations := &jsonObj.JsonObject{}
	if meta.HasKey("annotations") {
		annotations = meta.GetJsonObject("annotations")
	}
	annotations.Put(AnnAssumeTime, strconv.FormatInt(time.Now().UnixNano(), 10))
	annotations.Put(AnnAssignedFlag, "false")
	annotations.Put(ResourceUUID, result.GpuUuid[0])
	meta.Put("annotations", annotations.ToInterface())
	pod.Put("metadata", meta.ToInterface())

	_, err = decider.Client.UpdateResource(pod.ToString())
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
	gpuBytes, err := decider.Client.GetResource("GPU", GPUNamespace, gpuName)
	if err != nil {
		log.Fatalf("Failed to get GPU CRD, %s.", err)
	}
	gpu := kubesys.ToJsonObject(gpuBytes)
	status := gpu.GetJsonObject("status")
	allocated := status.GetJsonObject("allocated")

	oldMemoryStr, err := allocated.GetString("memory")
	if err != nil {
		log.Fatalf("Failed to get old memory, %s.", err)
	}
	oldCoreStr, err := allocated.GetString("core")
	if err != nil {
		log.Fatalf("Failed to get old core, %s.", err)
	}

	oldMemory, _ := strconv.ParseInt(oldMemoryStr, 10, 64)
	oldCore, _ := strconv.ParseInt(oldCoreStr, 10, 64)

	allocated.Put("memory", strconv.FormatInt(oldMemory + requestMemory, 10))
	allocated.Put("core", strconv.FormatInt(oldCore + requestCore, 10))

	status.Put("allocated", allocated.ToInterface())
	gpu.Put("status", status.ToInterface())

	_, err = decider.Client.UpdateResource(gpu.ToString())
	if err != nil {
		log.Fatalf("Failed to update GPU CRD, %s.", err)
	}

	log.Infof("Pod %s on namespace %s will run on node %s with %d gpu(s).", podName, namespace, result.NodeName, len(result.GpuUuid))

}

func (decider *Decider) deletePod(pod *jsonObj.JsonObject) {
	spec := pod.GetJsonObject("spec")
	schedulerName, err := spec.GetString("schedulerName")
	if err != nil || schedulerName != SchedulerName {
		return
	}
	nodeName, err := spec.GetString("nodeName")
	if err != nil {
		return
	}

	meta := pod.GetJsonObject("metadata")
	podName, err := meta.GetString("name")
	if err != nil {
		log.Fatalln("Failed to get pod name.")
	}
	namespace, err := meta.GetString("namespace")
	if err != nil {
		log.Fatalln("Failed to get pod namespace.")
	}

	annotations := meta.GetJsonObject("annotations")
	gpuUuid, err := annotations.GetString(ResourceUUID)
	if err != nil {
		log.Fatalf("Failed to get gpu uuid for pod %s on ns %s, %s.", podName, namespace, err)
	}

	requestMemory, requestCore := int64(0), int64(0)
	containers := spec.GetJsonArray("containers")
	for _, c := range containers.Values() {
		container := c.JsonObject()
		if !container.HasKey("resources") {
			continue
		}
		resources := container.GetJsonObject("resources")
		if !resources.HasKey("limits") {
			continue
		}
		limits := resources.GetJsonObject("limits")
		if val, err := limits.GetString(ResourceMemory); err == nil {
			m, _ := strconv.ParseInt(val, 10, 64)
			requestMemory += m
		}
		if val, err := limits.GetString(ResourceCore); err == nil {
			m, _ := strconv.ParseInt(val, 10, 64)
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
	gpuBytes, err := decider.Client.GetResource("GPU", GPUNamespace, gpuName)
	if err != nil {
		log.Fatalf("Failed to get GPU CRD, %s.", err)
	}

	gpu := kubesys.ToJsonObject(gpuBytes)
	status := gpu.GetJsonObject("status")
	allocated := status.GetJsonObject("allocated")

	oldMemoryStr, err := allocated.GetString("memory")
	if err != nil {
		log.Fatalf("Failed to get old memory, %s.", err)
	}
	oldCoreStr, err := allocated.GetString("core")
	if err != nil {
		log.Fatalf("Failed to get old core, %s.", err)
	}

	oldMemory, _ := strconv.ParseInt(oldMemoryStr, 10, 64)
	oldCore, _ := strconv.ParseInt(oldCoreStr, 10, 64)

	allocated.Put("memory", strconv.FormatInt(oldMemory - requestMemory, 10))
	allocated.Put("core", strconv.FormatInt(oldCore - requestCore, 10))

	status.Put("allocated", allocated.ToInterface())
	gpu.Put("status", status.ToInterface())

	_, err = decider.Client.UpdateResource(gpu.ToString())
	if err != nil {
		log.Fatalf("Failed to update GPU CRD, %s.", err)
	}

	log.Infof("Pod %s on namespace %s is deleled on node %s.", podName, namespace, nodeName)
}

func (decider *Decider) addGpu(gpu *jsonObj.JsonObject) {
	meta := gpu.GetJsonObject("metadata")
	gpuName, err := meta.GetString("name")
	if err != nil {
		log.Fatalln("Failed to get gpu name.")
	}

	spec := gpu.GetJsonObject("spec")
	gpuUuid, err := spec.GetString("uuid")
	if err != nil {
		log.Fatalf("Failed to get gpu %s's uuid, %s.", gpuName, err)
	}
	nodeName, err := spec.GetString("node")
	if err != nil {
		log.Fatalf("Failed to get gpu %s's node name, %s.", gpuName, err)
	}

	status := gpu.GetJsonObject("status")
	capacity := spec.GetJsonObject("capacity")
	allocated := status.GetJsonObject("allocated")

	coreCapacityStr, err := capacity.GetString("core")
	if err != nil {
		log.Fatalf("Failed to get gpu %s's core capacity, %s.", gpuName, err)
	}
	coreAllocatedStr, err := allocated.GetString("core")
	if err != nil {
		log.Fatalf("Failed to get gpu %s's core allocated, %s.", gpuName, err)
	}
	memoryCapacityStr, err := capacity.GetString("memory")
	if err != nil {
		log.Fatalf("Failed to get gpu %s's memory capacity, %s.", gpuName, err)
	}
	memoryAllocatedStr, err := allocated.GetString("memory")
	if err != nil {
		log.Fatalf("Failed to get gpu %s's memory allocated, %s.", gpuName, err)
	}

	coreCapacity, _ := strconv.ParseInt(coreCapacityStr, 10, 64)
	coreAllocated, _ := strconv.ParseInt(coreAllocatedStr, 10, 64)
	memoryCapacity, _ := strconv.ParseInt(memoryCapacityStr, 10, 64)
	memoryAllocated, _ := strconv.ParseInt(memoryAllocatedStr, 10, 64)

	hasDevicePlugin := false
	nodeBytes, err := decider.Client.GetResource("Node", "", nodeName)
	if err != nil {
		log.Fatalf("Failed to get GPU's node, %s.", err)
	}

	node := kubesys.ToJsonObject(nodeBytes)

	nodeStatus := node.GetJsonObject("status")
	nodeCapacity := nodeStatus.GetJsonObject("capacity")
	if nodeCapacity.HasKey(ResourceCore) {
		val, _ := nodeCapacity.GetString(ResourceCore)
		if val != "0" {
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

func (decider *Decider) modifyNode(node *jsonObj.JsonObject) {
	meta := node.GetJsonObject("metadata")
	nodeName, err := meta.GetString("name")
	if err != nil {
		log.Fatalln("Failed to get node name.")
	}

	hasDevicePlugin := false
	nodeStatus := node.GetJsonObject("status")
	nodeCapacity := nodeStatus.GetJsonObject("capacity")

	if nodeCapacity.HasKey(ResourceCore) {
		val, _ := nodeCapacity.GetString(ResourceCore)
		if val != "0" {
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
