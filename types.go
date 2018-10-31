package main

import (
	"encoding/json"
	"strconv"
)

//stub
type Result struct {
	Buffer string
}

//end stub

type ApiObjectSet struct {
	Kind       string       `json:"kind"`
	ApiObjects []*ApiObject `json:"items"`
}

type ApiObject struct {
	Kind     string    `json:"kind"`
	Metadata *Metadata `json:"metadata"`
}

type Metadata struct {
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
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
