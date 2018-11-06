package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file

const (
	CanaryStateBootstrap  string = "bootstrap"
	CanaryStateInProgress        = "in_progress"
	CanaryStateSuccess           = "success"
	CanaryStateFailed            = "failed"
)

const (
	CanaryPhaseTest     = "testing"
	CanaryPhaseAnalysis = "analysis"
)

// CanarySpec defines the desired state of Canary
type CanarySpec struct {
	// Target is the service that should be used for testing
	// purposes
	Target string `json:"target,omitempty"`

	TestPhase     TestPhaseSpec     `json:"testPhase,omitempty"`
	AnalysisPhase AnalysisPhaseSpec `json:"analysisPhase,omitempty"`
}

// TestPhaseSpec defines the configuration used to test the target
type TestPhaseSpec struct {
	Image string   `json:"image,omitempty"`
	Cmd   []string `json:"cmd,omitempty"`
}

// AnalysisPhaseSpec defines a resource that will accept metrics for further
// analysis purposes
type AnalysisPhaseSpec struct {
	Image string   `json:"image,omitempty"`
	Cmd   []string `json:"cmd,omitempty"`
}

// CanaryStatus defines the observed state of Canary
type CanaryStatus struct {
	Phase string `json:"phase,omitempty"`
	State string `json:"state,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Canary is the Schema for the canaries API
// +k8s:openapi-gen=true
type Canary struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CanarySpec   `json:"spec,omitempty"`
	Status CanaryStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CanaryList contains a list of Canary
type CanaryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Canary `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Canary{}, &CanaryList{})
}
