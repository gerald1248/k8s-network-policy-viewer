package main

import (
)

const (
	FilterIngress uint8 = 1 << iota
	FilterEgress
)

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
func filterEdgeMap(edgeMap *map[string][]string, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string, networkPolicies *[]ApiObject, globalNamespaces *[]string, mode int) {
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
		if len(o.Spec.PodSelector.MatchLabels) == 0 && len(o.Spec.PodSelector.MatchExpressions) == 0 {
			selectedPods = (*namespacePodMap)[namespace]
		} else if len(o.Spec.PodSelector.MatchExpressions) > 0 {
			// TODO: matchExpressions
			// TODO: IP range
		} else if len(o.Spec.PodSelector.MatchLabels) > 0 {
			selectedPods = selectPods(namespace, &o.Spec.PodSelector.MatchLabels, namespacePodMap, podLabelMap)
		}

		// source pods for ingress (.spec.from)

		// prepare set
		podsSet := make(map[string]struct{})
		for _, pod := range selectedPods {
			podsSet[pod] = struct{}{}
		}

		switch mode {
		case FilterIsolation: //first pass: isolation
			if flags&FilterIngress != 0 {
				isolated := len(o.Spec.Ingress) == 0 ||
					o.Spec.Ingress == nil
				// TODO: check for nested empty arrays
				if isolated {
					filterIngress(&podsSet, edgeMap)
				}
			}
			if flags&FilterEgress != 0 {
				isolated := len(o.Spec.Egress) == 0 ||
					o.Spec.Egress == nil
				if isolated {
					filterEgress(&podsSet, edgeMap)
				}
			}
		case FilterWhitelist: //second pass: whitelisting
			if flags&FilterIngress != 0 {
				for _, _ = range o.Spec.Ingress {
					// TODO: assume current namespace for now
					// TODO: ignore ports for now

					// special case: empty struct
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
				}
			}
			// TODO: egress whitelist
		}
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

// support among SDN providers still patchy?
// examine desired not necessarily actual state
// TODO: as with ingress, apply filter only once
func filterEgress(podsSet *map[string]struct{}, edgeMap *map[string][]string) {
	for pod, _ := range *podsSet {
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
