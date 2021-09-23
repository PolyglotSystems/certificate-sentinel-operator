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

// CertificateSentinelSpec defines the desired state of CertificateSentinel
type CertificateSentinelSpec struct {
	// Targets is the definition of K8s Objects to watch on the cluster and with what ServiceAccount
	Target Target `json:"target"`

	// Alerts is where the alerts will be sent to
	Alert Alert `json:"alert"`

	// ScanningInterval is how frequently the controller scans the cluster for these targets - defaults to 60s
	ScanningInterval int `json:"scanningInterval,omitempty"`

	// LogLevel controls the verbosity of the  - defaults to 1
	LogLevel int `json:"logLevel,omitempty"`
}

// Target provide what sort of objects we're watching for, be that a ConfigMap or a Secret
type Target struct {
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
}

// CertificateSentinelStatus defines the observed state of CertificateSentinel
type CertificateSentinelStatus struct {
	// DiscoveredCertificates is the slice of CertificateInformation that list the total set of discovered certificates
	DiscoveredCertificates []CertificateInformation `json:"discoveredCertificates"`
	// ExpiringCertificates is the number of certificates that are expiring
	ExpiringCertificates int `json:"expiringCertificates,omitempty"`
	// LastReportSent is last time the report was sent out
	LastReportSent int64 `json:"lastReportSent,omitempty"`
}

// CertificateInformation provides the status structure of what certificates have been discovered on the cluster
type CertificateInformation struct {
	// Namespace provides what namespace the certificate object was found in
	Namespace string `json:"namespace"`
	// Name provides the name of the certificate object
	Name string `json:"name"`
	// Kind provides the kind of the certificate object
	Kind string `json:"kind"`
	// APIVersion corresponds to the target kind apiVersion, so v1 is all really
	APIVersion string `json:"apiVersion"`
	// DataKey is the key for the data structure found
	DataKey string `json:"dataKey"`
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

// CertificateSentinel scans Secrets and ConfigMaps for x509 Certificates that are expiring and reports on them
type CertificateSentinel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertificateSentinelSpec   `json:"spec,omitempty"`
	Status CertificateSentinelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CertificateSentinelList contains a slice of CertificateSentinel
type CertificateSentinelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CertificateSentinel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CertificateSentinel{}, &CertificateSentinelList{})
}
