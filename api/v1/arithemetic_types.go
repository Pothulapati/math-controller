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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ArithemeticSpec defines the desired state of Arithemetic
type ArithemeticSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Expression is an math expression of Arithemetic, which the user wants to solve.
	Expression string `json:"expression,omitempty"`
}

// ArithemeticStatus defines the observed state of Arithemetic
type ArithemeticStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Answer is the solution to the expression
	Answer string `json:"answer"`
}

// +kubebuilder:object:root=true

// Arithemetic is the Schema for the arithemetics API
type Arithemetic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArithemeticSpec   `json:"spec,omitempty"`
	Status ArithemeticStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ArithemeticList contains a list of Arithemetic
type ArithemeticList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Arithemetic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Arithemetic{}, &ArithemeticList{})
}
