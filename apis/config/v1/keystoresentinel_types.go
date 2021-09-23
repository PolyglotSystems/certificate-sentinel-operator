/*
Copyright 2021 Polyglot Systems.

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

// KeystoreSentinelSpec defines the desired state of KeystoreSentinel
type KeystoreSentinelSpec struct {
	// Target is the K8s Objects to watch on the cluster
	Target KeystoreTarget `json:"target"`

	// Alert is where the alerts will be sent to
	Alert Alert `json:"alert"`

	// ScanningInterval is how frequently the controller scans the cluster for these targets - defaults to 30s
	ScanningInterval int `json:"scanningInterval,omitempty"`

	// LogLevel controls the verbosity of the  - defaults to 1
	LogLevel int `json:"logLevel,omitempty"`
}

// KeystoreTarget provide what sort of objects we're watching for, be that a ConfigMap or a Secret
type KeystoreTarget struct {
	// TargetName is a simple DNS/k8s compliant name for identification purposes
	TargetName string `json:"name"`
	// Namespaces is the slice of namespaces to watch on the cluster - can be a single wildcard to watch all namespaces
	Namespaces []string `json:"namespaces"`
	// NamespaceLabels is an optional slice of key pair labels to target, which will limit the scope of the matched namespaces to only ones with those labels
	NamespaceLabels []LabelSelector `json:"namespaceLabels,omitempty"`
	// Kind can be either ConfigMap or Secret
	Kind string `json:"kind"`
	// APIVersion corresponds to the target kind apiVersion, so v1 is all really
	APIVersion string `json:"apiVersion"`
	// TargetLabels is an optional slice of key pair labels to target, which will limit the scope of the matched objects to only ones with those labels
	TargetLabels []LabelSelector `json:"targetLabels,omitempty"`
	// ServiceAccount is the ServiceAccount to use in order to scan the cluster - this allows for separate RBAC per targeted object
	ServiceAccount string `json:"serviceAccount"`
	// DaysOut is the slice of days out alerts should be triggered at.  Defaults to 30, 60, and 90
	DaysOut []int `json:"daysOut,omitempty"`
	// KeystorePassword corresponds to the source for the the KeystorePassword
	KeystorePassword KeystorePassword `json:"keystorePassword"`
}

// KeystorePassword provides the input for the Keystore Password
type KeystorePassword struct {
	// Type could be 'secret', 'labels', or 'plaintext'
	Type      string          `json:"type"`
	Plaintext string          `json:"plaintext,omitempty"`
	Secret    SecretReference `json:"secretRef,omitempty"`
	Labels    LabelReference  `json:"labelRef,omitempty"`
}

// LabelReference provides the internal Secret reference to unlock the Keystore
type LabelReference struct {
	LabelSelectors []LabelSelector `json:"labelSelectors"`
	Key            string          `json:"key"`
}

// SecretReference provides the internal Secret reference to unlock the Keystore
type SecretReference struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// KeystoreSentinelStatus defines the observed state of KeystoreSentinel
type KeystoreSentinelStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DiscoveredKeystoreCertificates []KeystoreInformation `json:"discoveredKeystoreCertificates"`
	// ExpiringCertificates is the number of certificates that are expiring
	ExpiringCertificates int `json:"expiringCertificates,omitempty"`
	// TotalKeystoresFound is the number of Keystores found in scope
	TotalKeystoresFound int `json:"totalKeystoresFound,omitempty"`
	// KeystoresAtRisk is the number of Keystores that have expiring certificates
	KeystoresAtRisk int `json:"keystoresAtRisk,omitempty"`
	// LastReportSent is the time the report has been sent out by this Operator controller and when
	LastReportSent int64 `json:"lastReportSent,omitempty"`
}

// KeystoreInformation provides the status structure of what keystores have certificates that have been discovered on the cluster
type KeystoreInformation struct {
	// Namespace provides what namespace the Keystore object was found in
	Namespace string `json:"namespace"`
	// Name provides the name of the Keystore object
	Name string `json:"name"`
	// Kind provides the kind of the Keystore object
	Kind string `json:"kind"`
	// APIVersion corresponds to the target kind apiVersion, so v1 is all really
	APIVersion string `json:"apiVersion"`
	// DataKey is the key for the data structure found
	DataKey string `json:"dataKey"`
	// KeystoreAlias is the key for the data structure found
	KeystoreAlias string `json:"keystoreAlias"`
	// Expiration is the expiration date in YYYY-MM-DD
	Expiration string `json:"expiration"`
	// Name provides the name of the certificate object
	CommonName string `json:"commonName"`
	// CertificateAuthorityCommonName provides the Common Name of the signing Certificate Authority
	CertificateAuthorityCommonName string `json:"certificateAuthorityCommonName"`
	// IsCertificateAuthority returns a bool if the certificate is a CA
	IsCertificateAuthority bool `json:"isCertificateAuthority"`
	// TriggeredDaysOut provides the slice of days out that triggered the watch
	TriggeredDaysOut []int `json:"triggeredDaysOut,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KeystoreSentinel scans Secrets and ConfigMaps for Java Keystores with certificates that are expiring and reports on them
type KeystoreSentinel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeystoreSentinelSpec   `json:"spec,omitempty"`
	Status KeystoreSentinelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeystoreSentinelList contains a list of KeystoreSentinel
type KeystoreSentinelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeystoreSentinel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeystoreSentinel{}, &KeystoreSentinelList{})
}
