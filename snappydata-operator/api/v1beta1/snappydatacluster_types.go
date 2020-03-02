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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SnappyDataClusterSpec defines the desired state of SnappyDataCluster
type SnappyDataClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Size of the SnappyData cluster
	Size ClusterSize `json:"size"`
	// +optional
	SuspendAfter string `json:"suspendafter"`
	// +optional
	Resume bool   `json:"resume"`
	Image  string `json:"image"`
	// +optional
	// +kubebuilder:validation:Enum=Always;IfNotPresent;Never
	ImagePullPolicy v1.PullPolicy `json:"imapgePullPolicy,omitempty"`
	// +optional
	Configuration ClusterConfiguration `json:"configuration,omitempty"`
}

// ClusterSize describes the size of the SnappyData cluster
// Only one of the following cluster sizes can be specified,
// the default one is small
// +kubebuilder:validation:Enum=Small;Medium;Large;Custom
type ClusterSize string

const (
	// SmallCluster launches a small SnappyData cluster
	SmallCluster ClusterSize = "Small"
	// MediumCluster launaches a medium SnappyData cluster
	MediumCluster ClusterSize = "Medium"
	// LargeCluster launches a large SnappyData cluster
	LargeCluster ClusterSize = "Large"
	// CustomCluster launches a custom sized SnappyData cluster
	CustomCluster ClusterSize = "Custom"
)

// ClusterConfiguration specifies the configuration for clusters
type ClusterConfiguration struct {
	Locators MemberConfiguration `json:"locators,omitempty"`
	Servers  MemberConfiguration `json:"servers,omitempty"`
	Leads    MemberConfiguration `json:"leads,omitempty"`
}

// MemberConfiguration contains configuration SnappyData members
type MemberConfiguration struct {
	Count                int32                      `json:"count,omitempty"`
	Conf                 string                     `json:"conf,omitempty"`
	VolumeClaimTemplates []v1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
	Resources            v1.ResourceRequirements    `json:"resources,omitempty"`
}

// SnappyDataClusterStatus defines the observed state of SnappyDataCluster
type SnappyDataClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State ClusterState `json:"state"`
}

// ClusterState describes the state of the SnappyData cluster
// Only one of the following cluster statuses can be specified,
// +kubebuilder:validation:Enum=Running;Suspended
type ClusterState string

const (
	// RunningCluster specifies that the cluster is online
	RunningCluster ClusterState = "Running"

	// SuspendedCluster specifies that the cluster components have been stopped
	SuspendedCluster ClusterState = "Suspended"
)

// SnappyDataCluster is the Schema for the snappydataclusters API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster-Size",type=string,JSONPath=".spec.size",description="Size of the deployed cluster"
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=".status.state",description="State of the cluster"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SnappyDataCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnappyDataClusterSpec   `json:"spec,omitempty"`
	Status SnappyDataClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SnappyDataClusterList contains a list of SnappyDataCluster
type SnappyDataClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SnappyDataCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SnappyDataCluster{}, &SnappyDataClusterList{})
}
