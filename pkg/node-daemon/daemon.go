/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 **/

package node_daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	jsonutil "github.com/kubesys/kubernetes-client-go/pkg/util"
	v1 "github.com/kubesys/kubernetes-scheduler/pkg/apis/doslab.io/v1"
	log "github.com/sirupsen/logrus"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strconv"
	"strings"
	"syscall"
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
	// f, err := os.Create(GemSchedulerIpPath)
	// if err != nil {
	// 	log.Fatalf("Failed to create file %s, %s", GemSchedulerIpPath, err)
	// }
	// f.WriteString(os.Getenv(EnvGemSchedulerIp) + "\n")
	// f.Sync()
	// f.Close()

	err := os.MkdirAll(GemSchedulerGPUConfigPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create fir %s, %s.", GemSchedulerGPUConfigPath, err)
	}
	err = os.MkdirAll(GemSchedulerGPUPodManagerPortPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create fir %s, %s.", GemSchedulerGPUPodManagerPortPath, err)
	}

	log.Infoln("Loading NVML...")
	if err := nvml.Init(); err != nil {
		log.Infof("Failed to initialize NVML: %s.", err)
		log.Infof("If this is a GPU node, did you set the docker default runtime to `nvidia`?")
		log.Fatalln("Failed to discover gpus.")
	}
	defer func() {
		log.Infof("Shutdown of NVML returned: %s.", nvml.Shutdown())
	}()

	n, err := nvml.GetDeviceCount()
	if err != nil {
		log.Fatalf("Failed to get device count, %s.", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Failed to get hostname.")
	}

	for index := uint(0); index < n; index++ {
		device, err := nvml.NewDevice(index)
		if err != nil {
			log.Fatalf("Failed to new device, %s.", err)
		}

		_, err = os.Create(GemSchedulerGPUConfigPath + device.UUID)
		if err != nil {
			log.Fatalf("Failed to create file %s, %s", GemSchedulerGPUConfigPath+device.UUID, err)
		}

		_, err = os.Create(GemSchedulerGPUPodManagerPortPath + device.UUID)
		if err != nil {
			log.Fatalf("Failed to create file %s, %s", GemSchedulerGPUPodManagerPortPath+device.UUID, err)
		}

		gpu := v1.GPU{
			TypeMeta: metav1.TypeMeta{
				Kind:       "GPU",
				APIVersion: GPUCRDAPIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-gpu-%d", hostname, index),
				Namespace: GPUCRDNamespace,
			},
			Spec: v1.GPUSpec{
				UUID:   device.UUID,
				Model:  *device.Model,
				Family: getArchFamily(*device.CudaComputeCapability.Major, *device.CudaComputeCapability.Minor),
				Capacity: v1.R{
					Core:   "100",
					Memory: strconv.Itoa(int(*device.Memory)),
				},
				Node: hostname,
			},
			Status: v1.GPUStatus{
				Allocated: v1.R{
					Core:   "0",
					Memory: "0",
				},
			},
		}
		jb, err := json.Marshal(gpu)
		if err != nil {
			log.Fatalf("Failed to marshal gpu struct, %s.", err)
		}
		_, err = daemon.Client.CreateResource(string(jb))
		if err != nil && err.Error() != "request status 201 Created" {
			log.Fatalf("Failed to create gpu %s, %s.", gpu.Name, err)
		}
	}

	for {
		if daemon.PodMgr.queueOfModified.Len() > 0 {
			daemon.PodMgr.muOfModify.Lock()
			pod := daemon.PodMgr.queueOfModified.Remove()
			daemon.PodMgr.muOfModify.Unlock()
			time.Sleep(5 * time.Millisecond)
			go daemon.modifyPod(pod)
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

func (daemon *NodeDaemon) modifyPod(pod *jsonutil.ObjectNode) {
	meta := pod.GetObjectNode("metadata")
	if meta.Object["annotations"] == nil {
		return
	}
	annotations := meta.GetObjectNode("annotations")
	if annotations.Object[AnnAssumeTime] == nil {
		return
	}

	var str []string

	podName := meta.GetString("name")
	namespace := meta.GetString("namespace")
	str = append(str, namespace+"/"+podName)

	spec := pod.GetObjectNode("spec")
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

	if requestCore != 0 {
		str = append(str, strconv.FormatFloat(float64(requestCore)/100, 'f', 6, 64))
		str = append(str, strconv.FormatFloat(float64(requestCore)/100, 'f', 6, 64))
	}
	if requestMemory != 0 {
		str = append(str, strconv.Itoa(1024*requestMemory))
	}
	str[len(str)-1] += "\n"

	gpu := annotations.GetString(ResourceUUID)

	// Update gem-gpu-config file
	err := daemon.updateFile(str, GemSchedulerGPUConfigPath, gpu)
	if err != nil {
		log.Fatalf("Failed to update gem-gpu-config file, %s.", err)
	}

}

func (daemon *NodeDaemon) deletePod(pod *jsonutil.ObjectNode) {

}

func (daemon *NodeDaemon) updateFile(str []string, dir, gpu string) error {
	fileName := dir + "/" + gpu
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		return err
	}

	lines := make(map[string][]string)
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		words := strings.Split(line, " ")
		if len(words) == 1 {
			continue
		}
		lines[words[0]] = words[1:]
	}
	lines[str[0]] = str[1:]

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(strconv.Itoa(len(lines)) + "\n")
	if err != nil {
		return err
	}
	for k, v := range lines {
		s := k
		for i := 0; i < len(v); i++ {
			s += " "
			s += v[i]
		}
		_, err := f.WriteString(s)
		if err != nil {
			return err
		}
	}

	f.Sync()

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	if err != nil {
		return err
	}

	log.Infof("Success to update file %s.", fileName)

	return nil
}

func getArchFamily(computeMajor, computeMinor int) string {
	switch computeMajor {
	case 1:
		return "Tesla"
	case 2:
		return "Fermi"
	case 3:
		return "Kepler"
	case 5:
		return "Maxwell"
	case 6:
		return "Pascal"
	case 7:
		if computeMinor < 5 {
			return "volta"
		}
		return "Turing"
	case 8:
		return "Ampere"
	}
	return "Unknown"
}
