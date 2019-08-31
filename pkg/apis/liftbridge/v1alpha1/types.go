// Copyright 2019 The Liftbridge Operator Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LiftbridgeCluster represents a Liftbridge cluster.
type LiftbridgeCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec represents the desired state of the Liftbridge cluster.
	Spec LiftbridgeClusterSpec `json:"spec"`
	// Status provides information about the current state of the Liftbridge cluster.
	Status LiftbridgeClusterStatus `json:"status"`
}

// LiftbridgeClusterSpec represents the desired state of the Liftbridge cluster.
type LiftbridgeClusterSpec struct {
	// LogLevel is the log level to use for the Liftbridge cluster.
	// +optional
	LogLevel *string `json:"logLevel,omitempty"`
	// NATS is used to configure how the Liftbridge cluster connects to NATS.
	NATS LiftbridgeClusterNATSSpec `json:"nats"`
	// Paused indicates whether synchronisation of the current LiftbridgeCluster resource is currently paused.
	// +optional
	Paused *bool `json:"paused,omitempty"`
	// Replicas is the desired number of replicas in the Liftbridge cluster.
	Replicas int32 `json:"replicas"`
	// Storage is used to configure storage for the Liftbridge cluster.
	Storage LiftbridgeClusterStorageSpec `json:"storage"`
	// TLS is used to configure TLS for the Liftbridge cluster.
	// +optional
	TLS *LiftbridgeClusterTLSSpec `json:"tls,omitempty"`
	// Version is the version of Liftbridge in use.
	Version string `json:"version"`
}

// LiftbridgeClusterNATSSpec is used to configure how a Liftbridge cluster connects to NATS.
type LiftbridgeClusterNATSSpec struct {
	// Servers is the list of NATS servers Liftbridge should connect to.
	Servers []string `json:"servers"`
}

// LiftbridgeClusterTLSSpec is used to configure TLS for a Liftbridge cluster.
type LiftbridgeClusterTLSSpec struct {
	// SecretName is the name of the secret containing the TLS certificate and private key to use.
	SecretName string `json:"secretName"`
}

// LiftbridgeClusterStorageSpec is used to configure storage for a Liftbridge cluster.
type LiftbridgeClusterStorageSpec struct {
	// SizeGB is the desired size for the persistent volumes used for storing data.
	SizeGB int32 `json:"sizeGB"`
	// StorageClassName is the name of the storage class to use for provisioning persistent volumes.
	StorageClassName *string `json:"storageClassName"`
}

// LiftbridgeClusterStatus represents the status of the Liftbridge cluster.
type LiftbridgeClusterStatus struct {
	// Replicas is the current number of replicas in the Liftbridge cluster.
	Replicas int32 `json:"replicas"`
	// Selector is the label selector used to target replicas in the Liftbridge cluster.
	Selector string `json:"selector"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LiftbridgeClusterList represents a list of Liftbridge clusters.
type LiftbridgeClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items is the list of Liftbridge clusters.
	Items []LiftbridgeCluster `json:"items"`
}
