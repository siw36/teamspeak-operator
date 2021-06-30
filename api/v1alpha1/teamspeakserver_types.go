/*
Copyright 2021.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AutomaticRollout defines the rollout maintenance window for applying changes
type AutomaticRollout struct {
	Enabled                bool        `json:"enabled" default:"false"`
	IgnoreConnectedClients bool        `json:"ignoreConnectedClients" default:"false"`
	WindowStart            metav1.Time `json:"windowStart"`
	WindowStop             metav1.Time `json:"windowStop"`
}

// Persistence defines the persistent storage configuration
type Persistence struct {
	StorageClassName string                      `json:"storageClassName"`
	Resources        corev1.ResourceRequirements `json:"resources"`
	VolumeName       string                      `json:"volumeName,omitempty"`
	//AccessModes      corev1.[]PersistentVolumeAccessMode `json:"accessModes"`
}

// TeamSpeak defines the server and operator settings regarding a TeamSpeakServer
type TeamSpeak struct {
	Image            string            `json:"image" default:"docker.io/library/teamspeak:latest"`
	ImagePullPolicy  corev1.PullPolicy `json:"imagePullPolicy"`
	Persistence      Persistence       `json:"persistence"`
	AutomaticRollout AutomaticRollout  `json:"automaticRollout"`
}

// DatabaseConfigManaged defines the database that will be deployed and managed by the operator
type DatabaseConfigManaged struct {
	Image            string                      `json:"image" default:"docker.io/library/teamspeak:latest"`
	ImagePullPolicy  corev1.PullPolicy           `json:"imagePullPolicy"`
	Persistence      Persistence                 `json:"persistence"`
	StorageClassName string                      `json:"storageClassName"`
	Resources        corev1.ResourceRequirements `json:"resources"`
	VolumeName       string                      `json:"volumeName,omitempty"`
}

// DatabaseConfigUnmanaged defines the existing database
type DatabaseConfigUnmanaged struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	User               string `json:"user"`
	PasswordSecretName string `json:"passwordSecretName"`
	Database           string `json:"database"`
}

// Database defines the database backend for the TeamSpeakServer
type Database struct {
	Managed         bool                    `json:"managed" default:"true"`
	ConfigManaged   DatabaseConfigManaged   `json:"configManaged"`
	ConfigUnmanaged DatabaseConfigUnmanaged `json:"configUnmanaged"`
}

// TeamSpeakServerSpec defines the desired state of TeamSpeakServer
type TeamSpeakServerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	TeamSpeak TeamSpeak `json:"teamspeak"`
	Database  Database  `json:"database"`
}

// Defines the TeamSpeak deployment status
type Status struct {
	PodName     string      `json:"podName"`
	PodReady    bool        `json:"podReady"`
	Version     string      `json:"version"`
	LastRollout metav1.Time `json:"lastRollout"`
}

// DatabaseStatus defines databse status
type DatabaseStatus struct {
	PodName  string `json:"podName"`
	PodReady bool   `json:"podReady"`
}

// TeamSpeakServerStatus defines the observed state of TeamSpeakServer
type TeamSpeakServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	TeamSpeak  Status         `json:"teamSpeak"`
	Database   DatabaseStatus `json:"database"`
	LastUpdate metav1.Time    `json:"lastUpdated"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TeamSpeakServer is the Schema for the teamspeakservers API
type TeamSpeakServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamSpeakServerSpec   `json:"spec,omitempty"`
	Status TeamSpeakServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TeamSpeakServerList contains a list of TeamSpeakServer
type TeamSpeakServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TeamSpeakServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TeamSpeakServer{}, &TeamSpeakServerList{})
}
