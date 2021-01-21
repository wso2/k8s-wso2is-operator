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
// UserstoreSpec defines the desired state of Userstore
type UserstoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	TypeId      string                `json:"typeId"`
	Description string                `json:"description"`
	Name        string                `json:"name"`
	Properties  []UserstoreProperties `json:"properties"`
}

type Auth struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserstoreProperties struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// UserstoreStatus defines the observed state of Userstore
type UserstoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Userstore is the Schema for the userstores API
type Userstore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserstoreSpec   `json:"spec,omitempty"`
	Status UserstoreStatus `json:"status,omitempty"`

	Auth Auth `json:"auth"`
}

// +kubebuilder:object:root=true

// UserstoreList contains a list of Userstore
type UserstoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Userstore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Userstore{}, &UserstoreList{})
}
