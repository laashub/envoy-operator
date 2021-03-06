package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultContainerImage = "soloio/envoy:v0.1.6-131"

const (
	TLSCA           = "ca.crt"
	TLSCert         = "tls.cert"
	TLSKey          = "tls.key"
	EnvoyTLSVolPath = "/etc/certs/"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EnvoyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Envoy `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Envoy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              EnvoySpec   `json:"spec"`
	Status            EnvoyStatus `json:"status,omitempty"`
}

type EnvoySpec struct {
	ADSServer string `json:"adsServer"`
	ADSPort   int32  `json:"adsPort"`

	Image        string   `json:"image"`
	ImageCommand []string `json:"imageCommand"`

	// Secret name, containing ca cert, and potentially client cert and key with the names
	// ca.crt tls.crt, tls.key
	TLSSecretName string `json:"tls_secret_name"`

	AdminPort int32 `json:"adminPort"`

	ClusterIdTemplate string `json:"clusterIdTemplate"`

	NodeIdTemplate string `json:"nodeIdTemplate"`

	// Ports to expose on the service
	// If empty, no service will created for the Envoy
	// folllows format name: portnumber
	ServicePorts map[string]int32 `json:"servicePorts"`

	// StatsdSink string

	// OpenTracing string

	Deployment *EnvoyDeploymentSpec `json:"deployment,omitempty"`
	Injection  *InjectionSpec       `json:"ingress,omitempty"`
}

type EnvoyDeploymentSpec struct {
	// How many replicas of envoy we should have?
	Replicas uint32 `json:"replicas"`
}

// TODO: this is not implemented yet, but is written out to allow comments and discussion
type InjectionSpec struct {
	// This is should have configuration for how to inject.
	// for example:

	Mode           string // is the list below a whitelist or blacklist
	Namespaceslist []string
	// annotation per pod \ namespace that overrides above
	Annotation string
}

type EnvoyStatus struct {
}

// SetDefaults sets the default vaules for the Envoy spec and returns true if the spec was changed
func (e *Envoy) SetDefaults() bool {
	changed := false
	es := &e.Spec

	if es.Image == "" {
		es.Image = defaultContainerImage
		changed = true
	}
	if len(es.ImageCommand) == 0 {
		es.ImageCommand = []string{"/usr/local/bin/envoy"}
		changed = true
	}
	if es.AdminPort == 0 {
		es.AdminPort = 19000
		changed = true
	}
	if es.Injection == nil {
		if es.Deployment == nil {
			es.Deployment = &EnvoyDeploymentSpec{}
			changed = true
		}
		if es.Deployment.Replicas == 0 {
			es.Deployment.Replicas = 1
			changed = true
		}
	}
	return changed
}
