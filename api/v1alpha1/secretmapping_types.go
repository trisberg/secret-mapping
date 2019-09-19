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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SecretMappingSpec defines the desired state of SecretMapping
type SecretMappingSpec struct {
	// BindingSecret is the name of the Binding Secret that is created with the credentials
	BindingSecret string `json:"bindingSecret,omitempty"`
	// BindingPrefix is the prefix that will be prepended to the credentials properties
	BindingPrefix string `json:"bindingPrefix,omitempty"`
	// ServiceType is the type of the service, like 'mysql' or 'mongodb' etc.
	ServiceType string `json:"serviceType,omitempty"`
	// ServiceTypeKey is the key for the service type in the secret specified by SecretRef
	ServiceTypeKey string `json:"serviceTypeKey,omitempty"`
	// ServiceInstance is the instance name  of the service
	ServiceInstance string `json:"serviceInstance,omitempty"`
	// ServiceInstanceKey is the key for the service instance name in the secret specified by SecretRef
	ServiceInstanceKey string `json:"serviceInstanceKey,omitempty"`
	// SecretRef is a reference to a Secret containing the credentials
	SecretRef string `json:"secretRef,omitempty"`
	// URI is the service URI that can be used to connect to the service
	URI string `json:"uri,omitempty"`
	// URIKey is the key for the URI in the secret specified by SecretRef
	URIKey string `json:"uriKey,omitempty"`
	// PasswordKey is the key for the password in the secret specified by SecretRef
	PasswordKey string `json:"passwordKey,omitempty"`
	// Username is the username to use for connecting to the service
	Username string `json:"username,omitempty"`
	// UsernameKey is the key for the username in the secret specified by SecretRef
	UsernameKey string `json:"usernameKey,omitempty"`
	// Host is the hostname or IP address for the service
	Host string `json:"host,omitempty"`
	// HostKey is the key for the host in the secret specified by SecretRef
	HostKey string `json:"hostKey,omitempty"`
	// Port is the port used by the service
	Port int `json:"port,omitempty"`
	// PortKey is the key for the port in the secret specified by SecretRef
	PortKey string `json:"portKey,omitempty"`
	// Important: Run "make" to regenerate code after modifying this file
}

// SecretMappingStatus defines the observed state of SecretMapping
type SecretMappingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// SecretMapping is the Schema for the secretmappings API
type SecretMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretMappingSpec   `json:"spec,omitempty"`
	Status SecretMappingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecretMappingList contains a list of SecretMapping
type SecretMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretMapping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretMapping{}, &SecretMappingList{})
}
