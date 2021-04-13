package v1
import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TaskSpec defines the desired state of Task
type TaskSpec struct {
}

// TaskStatus defines the observed state of Task
type TaskStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	PodStatus      v1.PodStatus `json:"pod_status,omitempty"`
	GPUId          string       `json:"gpu_id,omitempty"`
	PodManagerPort int          `json:"port,omitempty"`
}

// +kubebuilder:object:root=true

// Task is the Schema for the tasks API
type Task struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   v1.PodSpec `json:"spec,omitempty"`
	Status TaskStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskList contains a list of Task
type TaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Task `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Task{}, &TaskList{})
}
