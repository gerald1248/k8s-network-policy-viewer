package main

import (
	"strings"
)

// enum for ingress/egress policy
const (
	FilterIngress uint8 = 1 << iota
	FilterEgress
)

// enum for black-/whitelisting
const (
	FilterIsolation int = 0
	FilterWhitelist     = 1
)

func initializeEdgeMap(edgeMap *map[string][]string, namespacePodMap *map[string][]string) {
	var allPods []string
	for _, v := range *namespacePodMap {
		for _, s := range v {
			allPods = append(allPods, s)
		}
	}
	for _, outer := range allPods {
		for _, inner := range allPods {
			if inner == outer {
				continue
			}
			(*edgeMap)[outer] = append((*edgeMap)[outer], inner)
		}
	}
}

// filterEdgeMap is called twice: isolation 1st, whitelisting 2nd
// see central switch
func filterEdgeMap(
	edgeMap *map[string][]string,
	namespacePodMap *map[string][]string,
	namespaceLabelMap *map[string]map[string]string,
	podLabelMap *map[string]map[string]string,
	networkPolicies *[]APIObject,
	mode int) {
	for _, o := range *networkPolicies {
		namespace := o.Metadata.Namespace

		// set policy state
		var flags uint8
		if len(o.Spec.PolicyTypes) == 0 {
			flags |= FilterIngress
		} else {
			for _, policyType := range o.Spec.PolicyTypes {
				switch policyType {
				case "Ingress":
					flags |= FilterIngress
				case "Egress":
					flags |= FilterEgress
				}
			}
		}

		// destination pods in current namespace
		// .spec.podSelector is mandatory, so assume it's present
		var selectedPods []string
		// TODO: parse matchExpressions
		if len(o.Spec.PodSelector.MatchLabels) == 0 {
			selectedPods = (*namespacePodMap)[namespace]
		} else if len(o.Spec.PodSelector.MatchLabels) > 0 {
			selectedPods = selectPods(namespace, o.Spec.PodSelector, namespacePodMap, podLabelMap)
		}

		// source pods for ingress (.spec.from)

		// prepare set
		podsSet := make(map[string]struct{})
		for _, pod := range selectedPods {
			podsSet[pod] = struct{}{}
		}

		switch mode {
		case FilterIsolation: //first pass: isolation
			// TODO: check if policyTypes property always up-to-date
			if flags&FilterIngress != 0 {
				filterIngress(&podsSet, edgeMap)
			}
			if flags&FilterEgress != 0 {
				filterEgress(&podsSet, edgeMap)
			}
		case FilterWhitelist: //second pass: whitelisting
			// TODO: ignore ports for now
			if flags&FilterIngress != 0 {
				if o.Spec.Ingress != nil && len(o.Spec.Ingress) > 0 {
					for _, rule := range o.Spec.Ingress {
						if len(rule.From) == 0 {
							// TODO: assume current namespace for now
							namespacePods := (*namespacePodMap)[namespace]
							for _, namespacePod := range namespacePods {
								for _, selectedPod := range selectedPods {
									if namespacePod == selectedPod {
										continue
									}
									(*edgeMap)[namespacePod] = append((*edgeMap)[namespacePod], selectedPod)
								}
								(*edgeMap)[namespacePod] = unique((*edgeMap)[namespacePod])
							}
						} else {
							// identify source pods
							for _, peer := range rule.From {
								var fromPods []string
								if peer.NamespaceSelector == nil || peer.NamespaceSelector.MatchLabels == nil {
									fromPods = selectPods(namespace, peer.PodSelector, namespacePodMap, podLabelMap)
								} else {
									namespaces := selectNamespaces(&peer.NamespaceSelector.MatchLabels, namespacePodMap, namespaceLabelMap)
									// v1.10 does not support mixed namespace/pod selection
									// TODO: don't use pointer; pass in selector
									var selectorStub Selector
									fromPods = selectPodsAcrossNamespaces(&namespaces, &selectorStub, namespacePodMap, podLabelMap)
								}
								for _, fromPod := range fromPods {
									for _, selectedPod := range selectedPods {
										if fromPod == selectedPod {
											continue
										}
										(*edgeMap)[fromPod] = append((*edgeMap)[fromPod], selectedPod)
									}
									(*edgeMap)[fromPod] = unique((*edgeMap)[fromPod])
								}
							}
						}
					}
				}
			}
		}
		// TODO: egress whitelist
	}
}

// TODO: apply filter only once when multiple network policies apply to one namespace
func filterIngress(podsSet *map[string]struct{}, edgeMap *map[string][]string) {
	for fromString, toSlice := range *edgeMap {
		arr := []string{}
		for _, pod := range toSlice {
			if _, ok := (*podsSet)[pod]; !ok {
				arr = append(arr, pod)
			}
		}
		(*edgeMap)[fromString] = arr
	}
}

// brute force deduplication
// TODO: refactor
func deduplicateEdgeMap(edgeMap *map[string][]string) {
	for k, v := range *edgeMap {
		(*edgeMap)[k] = unique(v)
	}
}

func filterIntraNamespace(edgeMap *map[string][]string) {
	var namespaceFrom, namespaceTo string
	var aFrom, aTo []string
	for k, v := range *edgeMap {
		var aNew []string
		aFrom = strings.Split(k, namespaceSeparator())
		namespaceFrom = aFrom[0]
		for _, item := range (*edgeMap)[k] {
			aTo = strings.Split(item, namespaceSeparator())
			namespaceTo = aTo[0]
			if namespaceFrom != namespaceTo {
				aNew = append(aNew, item)
			}
		}
		if len(aNew) != len(v) {
			(*edgeMap)[k] = aNew
		}
	}
}

// support among SDN providers still patchy?
// examine desired not necessarily actual state
// TODO: as with ingress, apply filter only once
func filterEgress(podsSet *map[string]struct{}, edgeMap *map[string][]string) {
	for pod := range *podsSet {
		(*edgeMap)[pod] = nil
	}
}

func unique(slice []string) []string {
	keys := make(map[string]struct{})
	list := []string{}
	for _, entry := range slice {
		if _, ok := keys[entry]; !ok {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
