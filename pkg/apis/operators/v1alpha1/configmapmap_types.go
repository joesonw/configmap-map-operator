package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigMapMapSpecItem defines item to be merged
type ConfigMapMapSpecItem struct {
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	SubPath   string `json:"subPath,omitempty"`
}

// ConfigMapMapSpec defines the desired state of ConfigMapMap
type ConfigMapMapSpec struct {
	Namespace string                          `json:"namespace,omitempty"`
	Name      string                          `json:"name,omitempty"`
	Data      map[string]ConfigMapMapSpecItem `json:"data,omitempty"`
}

// ConfigMapMapStatus defines the observed state of ConfigMapMap
type ConfigMapMapStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapMap is the Schema for the configmapmaps API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=configmapmaps,scope=Namespaced
type ConfigMapMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigMapMapSpec   `json:"spec,omitempty"`
	Status ConfigMapMapStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapMapList contains a list of ConfigMapMap
type ConfigMapMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigMapMap `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigMapMap{}, &ConfigMapMapList{})
}
