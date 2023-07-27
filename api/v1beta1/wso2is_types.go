/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Wso2IsSpec defines the desired state of Wso2Is
type Wso2IsSpec struct {
	Size int32 `json:"replicas"`

	Resources Resources `json:"resources"`
	Hpa       Hpa       `json:"hpa,omitempty"`

	// +kubebuilder:default:="6.1.0"
	Version        string              `json:"version,omitempty"`
	Configurations Configurations      `json:"configurations"`
	TomlConfig     string              `json:"tomlConfig,omitempty"`
	TomlConfigFile string              `json:"tomlConfigFile,omitempty"`
	KeystoreMounts []KeystoreMount     `json:"keystoreMounts,omitempty"`
	Template       TemplateAnnotations `json:"template,omitempty"`
}

type Resources struct {
	Limits   ResourceLimits   `json:"limits"`
	Requests ResourceRequests `json:"requests"`
}

type ResourceLimits struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type ResourceRequests struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

type Hpa struct {
	MinReplicas              int32 `json:"minReplicas"`
	MaxReplicas              int32 `json:"maxReplicas"`
	CpuUtilizationPercentage int32 `json:"cpuUtilizationPercentage"`
}

type TemplateAnnotations struct {
	Annotations map[string]string `json:"annotations,omitempty"`
}

type KeystoreMount struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Wso2IsStatus defines the observed state of Wso2Is
type Wso2IsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes           []string `json:"nodes" toml:"nodes"`
	ServiceName     string   `json:"serviceName"`
	IngressHostname string   `json:"ingressHostname"`
	Replicas        string   `json:"replicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Service Name",type="string",JSONPath=`.status.serviceName`
// +kubebuilder:printcolumn:name="Ingress Hostname",type="string",JSONPath=`.status.ingressHostname`
// +kubebuilder:printcolumn:name="Desired Replicas",type="string",JSONPath=`.spec.replicas`
// +kubebuilder:printcolumn:name="Current Replicas",type="string",JSONPath=`.status.replicas`
// +kubebuilder:printcolumn:name="Host Name",type="string",JSONPath=`.spec.configurations.host`

// Wso2Is is the Schema for the wso2is API
type Wso2Is struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Wso2IsSpec   `json:"spec,omitempty"`
	Status Wso2IsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// Wso2IsList contains a list of Wso2Is
type Wso2IsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Wso2Is `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Wso2Is{}, &Wso2IsList{})
}
