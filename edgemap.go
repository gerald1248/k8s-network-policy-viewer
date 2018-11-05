package main

import ()

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

func filterEdgeMap(edgeMap *map[string][]string, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string, networkPolicies *[]ApiObject) {
	for _, o := range *networkPolicies {
		podsSet := map[string]bool{}
		namespace := o.Metadata.Namespace
		for _, pod := range (*namespacePodMap)[namespace] {
			podsSet[pod] = true
		}
		// 1. apply blanket ingress/egress policies
		if len(o.Spec.PolicyTypes) == 0 {
			// if none specified, default to ingress policy
			filterIngress(&podsSet, edgeMap)
		} else {
			for _, policyType := range o.Spec.PolicyTypes {
				switch policyType {
				case "Ingress":
					filterIngress(&podsSet, edgeMap)
				case "Egress":
					filterEgress(&podsSet, edgeMap)
				}
			}
		}
		// 2. now deal with whitelisted pods
		//selectedPods := selectPods(namespace, &o.Spec.PodSelector.MatchLabels, namespacePodMap, podLabelMap)
	}
}

// TODO: apply filter only once when a namespace follows multiple network policies
func filterIngress(podsSet *map[string]bool, edgeMap *map[string][]string) {
	for k, v := range *edgeMap {
		for i, pod := range v {
			// fast lookup at the cost of a resource-intensive data structure
			if (*podsSet)[pod] {
				(*edgeMap)[k] = append((*edgeMap)[k][:i], (*edgeMap)[k][:i+1]...)
			}
		}
	}
}

// support among SDN providers still patchy?
// examine desired not necessarily actual state
// TODO: as with ingress, apply filter only once
func filterEgress(podsSet *map[string]bool, edgeMap *map[string][]string) {
	for pod, _ := range *podsSet {
		(*edgeMap)[pod] = nil
	}
}
