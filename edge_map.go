package main

import (
)

const (
	FilterIngress uint8 = 1 << iota
	FilterEgress
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

func filterEdgeMap(edgeMap *map[string][]string, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string, networkPolicies *[]ApiObject, globalNamespaces *[]string) {
	ingressPods := []string{}
	egressPods := []string{}

	// handle global namespaces
	for _, globalNamespace := range *globalNamespaces {
		ingressPods = append(ingressPods, (*namespacePodMap)[globalNamespace]...)
		egressPods = append(egressPods, (*namespacePodMap)[globalNamespace]...)
	}

	for _, o := range *networkPolicies {
		var flags uint8
		podsSet := make(map[string]struct{})
		namespace := o.Metadata.Namespace
		for _, pod := range (*namespacePodMap)[namespace] {
			podsSet[pod] = struct{}{}
		}
		// 1. apply blanket ingress/egress policies
		flags = 0
		if len(o.Spec.PolicyTypes) == 0 {
			// if none specified, default to ingress policy
			filterIngress(&podsSet, edgeMap)
			flags |= FilterIngress
		} else {
			for _, policyType := range o.Spec.PolicyTypes {
				switch policyType {
				case "Ingress":
					filterIngress(&podsSet, edgeMap)
					flags |= FilterIngress
				case "Egress":
					filterEgress(&podsSet, edgeMap)
					flags |= FilterEgress
				}
			}
		}
		// consider only namespace-wide filters for now
		if flags&FilterIngress == 0 {
			ingressPods = unique(append(ingressPods, (*namespacePodMap)[namespace]...))
		}
		if flags&FilterEgress == 0 {
			egressPods = append(egressPods, (*namespacePodMap)[namespace]...)
		}

		// 2. now deal with whitelisted pods
			selectedPods := selectPods(namespace, &o.Spec.PodSelector.MatchLabels, namespacePodMap, podLabelMap)
			for _, pod := range selectedPods {
				if len(o.Spec.Ingress) > 0 {
					// empty ingress definition: all pods in all namespaces
					if o.Spec.Ingress[0].From == nil {
						for _, egressPod := range egressPods {
							if egressPod == pod {
								continue
							}
							slice := (*edgeMap)[egressPod]
							slice = unique(append(slice, pod))
							(*edgeMap)[egressPod] = slice
						}
					}
					// TODO: ingress selector
				}
				// TODO: egress
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
