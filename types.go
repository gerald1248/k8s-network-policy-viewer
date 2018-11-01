package main

import (
	"encoding/json"
	"strconv"
)

// nested custom type
// union of relevant fields for objects of kinds
// Pod and NetworkPolicy
type ApiObjectSet struct {
	Kind       string       `json:"kind"`
	ApiObjects []*ApiObject `json:"items"`
}

type ApiObject struct {
	Kind     string    `json:"kind"`
	Metadata *Metadata `json:"metadata"`
	Spec     *Spec     `json:"spec"`
	Status   *Status   `json:"status"`
}

type Metadata struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
}

type Status struct {
	ContainerStatuses []*ContainerStatus `json:"containerStatuses"`
}

type ContainerStatus struct {
	Ready bool `json:"ready"`
}

// the Spec struct is where the union of various 'kinds' of API object
// is most apparent; Ingress features most prominently
type Spec struct {
	PodSelector *Selector `json:"podSelector",omitempty`
	PolicyTypes []string  `json:"policyTypes",omitempty`
	Ingress []*Ingress    `json:"ingress",omitempty`
	Egress []*Egress      `json:"egress",omitempty`
}

type Selector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}

type Ingress struct {
	From []*SelectorCollection `json:"from"`
}

type Egress struct {
	To []*SelectorCollection `json:"to"`
}

type SelectorCollection struct {
	NamespaceSelector *Selector `json:"namespaceSelector"`
	PodSelector *Selector       `json:"podSelector"`
}

type MinimalObject struct {
	Kind string
}

type Table struct {
	Row []string
}

type CoerceString struct {
	s string
}

func (cs *CoerceString) String() string {
	return cs.s
}

//see also: kubernetes/api/util.go for fuzzy alternative
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

// TODO: remove?
type ContainerSet []ContainerSpec

type ContainerSpec struct {
	Namespace string
	Name      string
	Container string
}
