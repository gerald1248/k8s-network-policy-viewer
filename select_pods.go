package main

import (
	"log"
)

func selectPods(namespace string, selector *map[string]string, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string {
	// special case: empty map selects all pods
	if len(*selector) == 0 {
		return (*namespacePodMap)[namespace]
	}

	// alternatively, select pods matching at least one label pair
	var selectedPods []string
	for _, pod := range (*namespacePodMap)[namespace] {
		labels := (*podLabelMap)[pod]
		for k, v := range *selector {
			if len(k) > 0 && labels[k] == v {
				log.Printf("%s => %s\n", k, v)
				selectedPods = append(selectedPods, pod)
			}
		}
	}
	if len(selectedPods) > 0 {
		log.Printf("Selected pods: %v\n", selectedPods)
	}
	return selectedPods
}
