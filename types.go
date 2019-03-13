// ingress/egress rule types based on:
// https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/networking/types.go

package main

import (
	"encoding/json"
	"strconv"
)

// APIObjectSet is a nested custom type
// union of relevant fields for objects of kinds
// Pod and NetworkPolicy
type APIObjectSet struct {
	Kind       string       `json:"kind"`
	APIObjects []*APIObject `json:"items"`
}

// APIObject carries the fixed top-level properties
type APIObject struct {
	Kind     string    `json:"kind"`
	Metadata *Metadata `json:"metadata"`
	Spec     *Spec     `json:"spec"`
	Status   *Status   `json:"status"`
}

// Metadata wraps the essential metadata fields
type Metadata struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
}

// Status wraps status entries
type Status struct {
	ContainerStatuses []*ContainerStatus `json:"containerStatuses"`
}

// ContainerStatus here flags only 'ready' but could pick up warnings, memory pressure, etc.
type ContainerStatus struct {
	Ready bool `json:"ready"`
}

// Spec is where the union of various 'kinds' of API object
// is most apparent; Ingress features most prominently
type Spec struct {
	PodSelector *Selector                   `json:"podSelector"`
	PolicyTypes []string                    `json:"policyTypes"`
	Ingress     []*NetworkPolicyIngressRule `json:"ingress"`
	Egress      []*NetworkPolicyEgressRule  `json:"egress"`
}

// NetworkPolicyIngressRule is a core component of the network policy construct
type NetworkPolicyIngressRule struct {
	// TODO: Ports []NetworkPolicyPort
	From []*NetworkPolicyPeer `json:"from"`
}

// NetworkPolicyEgressRule is not expected to be used often
type NetworkPolicyEgressRule struct {
	// TODO: Ports []NetworkPolicyPort
	To []*NetworkPolicyPeer `json:"to"`
}

// NetworkPolicyPeer wraps the pod and namespace selectors
type NetworkPolicyPeer struct {
	PodSelector       *Selector `json:"podSelector"`
	NamespaceSelector *Selector `json:"namespaceSelector"`
	// TODO: IPBlock
}

// Selector is used in various selection contexts
type Selector struct {
	MatchLabels      map[string]string           `json:"matchLabels"`
	MatchExpressions []*LabelSelectorRequirement `json:"matchExpressions"`
}

// LabelSelectorRequirement wraps a key--operator--values selector
type LabelSelectorRequirement struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// Port is a leaf node that usually takes 8080/TCP
type Port struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// MinimalObject lacks everything except a Kind property for traversing sets of heterogeneous objects
type MinimalObject struct {
	Kind string
}

// Table is here limited to a series of rows
type Table struct {
	Row []string
}

// CoerceString takes only a string parameter
type CoerceString struct {
	s string
}

func (cs *CoerceString) String() string {
	return cs.s
}

//UnmarshalJSON - see also: kubernetes/api/util.go for fuzzy alternative
func (cs *CoerceString) UnmarshalJSON(value []byte) error {
	if value[0] == '"' {
		return json.Unmarshal(value, &cs.s)
	}
	var i int
	err := json.Unmarshal(value, &i)
	if err == nil {
		cs.s = strconv.Itoa(i)
		return nil
	}
	return err
}

// Result captures the metrics output of the policy scan
type Result struct {
	PercentageIsolated          int `json:"percentageIsolated"`
	PercentageNamespaceIsolated int `json:"percentageNamespaceIsolated"`
	PercentageNamespaceCoverage int `json:"percentageNamespaceCoverage"`
}
