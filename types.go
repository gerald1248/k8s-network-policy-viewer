package main

import (
	"encoding/json"
	"regexp"
	"strconv"
)

//stub
type Result struct {
        Namespace string
        Name      string
}
//end stub

type Metadata struct {
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

type Container struct {
	Image           string           `json:"image"`
	Name            string           `json:"name"`
	ImagePullPolicy string           `json:"imagePullPolicy"`
	Env             []*EnvItem       `json:"env"`
	Resources       *Resources       `json:"resources"`
	LivenessProbe   *Probe           `json:"livenessProbe"`
	ReadinessProbe  *Probe           `json:"readinessProbe"`
	SecurityContext *SecurityContext `json:"securityContext"`
}

type Probe struct {
	TimeoutSeconds      int `json:"timeoutSeconds"`
	PeriodSeconds       int `json:"periodSeconds"`
	SuccessThreshold    int `json:"successThreshold"`
	InitialDelaySeconds int `json:"initialDelaySeconds"`
	FailureThreshold    int `json:"failureThreshold"`
}

type SecurityContext struct {
	Privileged   bool  `json:"privileged"`
	RunAsNonRoot bool  `json:"runAsNonRoot"`
	RunAsUser    int64 `json:"runAsUser"`
}

type EnvItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ContainerSet []ContainerSpec

type ContainerSpec struct {
	Namespace string
	Name      string
	Container string
}

type Resources struct {
	Limits   *ResourceConstraint `json:"limits"`
	Requests *ResourceConstraint `json:"requests"`
}

type ResourceConstraint struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type Port struct {
	TargetPort CoerceString `json:"targetPort"`
}

func (r *ResourceConstraint) Complete() bool {
	return len(r.CPU) > 0 && len(r.Memory) > 0
}

func (r *ResourceConstraint) Valid() bool {
	//see Kubernetes pkg/api/validation/validation.go ll. 1300f.
	reCpu := regexp.MustCompile(`^[0-9]+m?$`)
	reMemory := regexp.MustCompile(`^[0-9]+(k|M|G|T|P|E|Ki|Mi|Gi|Ti|Pi|Ei)?$`)

	if len(r.CPU) > 0 && reCpu.FindStringIndex(r.CPU) == nil {
		return false
	}

	if len(r.Memory) > 0 && reMemory.FindStringIndex(r.Memory) == nil {
		return false
	}

	return true
}

func (r *ResourceConstraint) None() bool {
	return len(r.CPU) == 0 && len(r.Memory) == 0
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
