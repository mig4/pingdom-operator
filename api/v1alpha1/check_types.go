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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Type of check.
// +kubebuilder:validation:Enum=http;httpcustom;tcp;ping;dns;udp;smtp;pop3;imap
type CheckType string

const (
	// Check response to a HTTP request
	Http CheckType = "http"

	// Check response to a custom HTTP request
	HttpCustom CheckType = "httpcustom"

	// Send a packet to a TCP port
	Tcp CheckType = "tcp"

	// Send a ping (ICMP request) to the host
	Ping CheckType = "ping"

	// Try to resolve host using specified DNS server
	Dns CheckType = "dns"

	// Send a packet to a UDP port
	Udp CheckType = "udp"

	// Open a connection to SMTP server
	Smtp CheckType = "smtp"

	// Open a connection to a POP3 server
	Pop3 CheckType = "pop3"

	// Open a connection to an IMAP server
	Imap CheckType = "imap"
)

// CheckSpec defines the desired state of Check
type CheckSpec struct {
	// Check name; defaults to name of the object in Kubernetes
	// +optional
	Name *string `json:"name"`

	// Target host
	Host string `json:"host"`

	// Type of check, can be one of:
	// http, httpcustom, tcp, ping, dns, udp, smtp, pop3, imap
	Type CheckType `json:"type"`

	// Paused; defaults to false.
	// +optional
	Paused *bool `json:"paused,omitempty"`

	// Target port
	// Required for check types: tcp, udp
	// Optional for: http(80), httpcustom(80), smtp(25), pop3(110), imap(143)
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	Port *int32 `json:"port,omitempty"`

	// HTTP Checks

	// Target path on server
	// Defaults to `/`.
	// +optional
	Url *string `json:"url,omitempty"`

	// Connection encryption; defaults to false
	// +optional
	Encryption *bool `json:"encryption,omitempty"`

	// Secret storing Pingdom API credentials
	CredentialsSecret corev1.LocalObjectReference `json:"credentialsSecret"`
}

// Status/result of a check.
// +kubebuilder:validation:Enum=up;down;unconfirmed_down;unknown;paused
type CheckResult string

// CheckStatus defines the observed state of Check
type CheckStatus struct {
	// Check identifier
	Id int32 `json:"id"`

	// Check name
	// +optional
	Name *string `json:"name"`

	// Type of check
	Type CheckType `json:"type"`

	// Target host
	Host string `json:"host"`

	// Target port
	// +optional
	Port *int32 `json:"port,omitempty"`

	// Current check status
	Status CheckResult `json:"status"`

	// Paused; defaults to false.
	// +optional
	Paused *bool `json:"paused,omitempty"`

	// Timestamp of last error (if any).
	// +optional
	LastErrorTime *metav1.Time `json:"lasterrortime,omitempty"`

	// Timestamp of last test (if any).
	// +optional
	LastTestTime *metav1.Time `json:"lasttesttime,omitempty"`

	// Response time (in milliseconds) of last test.
	// +optional
	LastResponseTimeMilis *int64 `json:"lastresponsetime,omitempty"`

	// Check creation time.
	CreatedTime metav1.Time `json:"created"`

	// HTTP Checks

	// Target path on server
	// Defaults to `/`.
	// +optional
	Url *string `json:"url,omitempty"`

	// Connection encryption; defaults to false
	// +optional
	Encryption *bool `json:"encryption,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ID",type=string,JSONPath=`.status.id`,description="Check ID"
// +kubebuilder:printcolumn:name="type",type=string,JSONPath=`.status.type`,description="Check type"
// +kubebuilder:printcolumn:name="status",type=string,JSONPath=`.status.status`,description="Check status"
// +kubebuilder:printcolumn:name="paused",type=boolean,JSONPath=`.status.paused`,description="Pasused"

// Check is the Schema for the checks API
type Check struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CheckSpec   `json:"spec,omitempty"`
	Status CheckStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CheckList contains a list of Check
type CheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Check `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Check{}, &CheckList{})
}
