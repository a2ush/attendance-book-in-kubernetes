/*
Copyright 2022 a2ush.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AttendanceBookSpec defines the desired state of AttendanceBook
type AttendanceBookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AttendanceBook. Edit attendancebook_types.go to remove/update
	//+kubebuilder:validation:Enum=present;absent
	Attendance string `json:"attendance"`
	//+kubebuilder:default=BLANK
	// +optional
	Reason string `json:"reason"`
}

// AttendanceBookStatus defines the observed state of AttendanceBook
type AttendanceBookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Attendance string `json:"attendance"`
	// +optional
	Reason string `json:"reason"`
}

//+kubebuilder:resource:shortName=ab
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="ATTENDANCE",type="string",JSONPath=".status.attendance"
//+kubebuilder:printcolumn:name="REASON",type="string",JSONPath=".status.reason"
//+kubebuilder:printcolumn:name="REPORT TIME",type="string",JSONPath=".metadata.creationTimestamp"

// AttendanceBook is the Schema for the attendancebooks API
type AttendanceBook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AttendanceBookSpec   `json:"spec,omitempty"`
	Status AttendanceBookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AttendanceBookList contains a list of AttendanceBook
type AttendanceBookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AttendanceBook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AttendanceBook{}, &AttendanceBookList{})
}
