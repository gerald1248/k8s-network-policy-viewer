package main

func selectPods(namespace string, selector *map[string]string, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string {
	// special case: empty map selects all pods
	if *selector == nil || len(*selector) == 0 {
		return (*namespacePodMap)[namespace]
	}

	// alternatively, select pods matching at least one label pair
	var selectedPods []string
	for _, pod := range (*namespacePodMap)[namespace] {
		labels := (*podLabelMap)[pod]
		for k, v := range *selector {
			if len(k) > 0 && labels[k] == v {
				selectedPods = append(selectedPods, pod)
			}
		}
	}
	return selectedPods
}

func selectPodsAcrossNamespaces(namespaces *[]string, selector *map[string]string, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string {
	var allPods []string
	for _, namespace := range *namespaces {
		selectedPods := selectPods(namespace, selector, namespacePodMap, podLabelMap)
		allPods = append(allPods, selectedPods...)
	}
	return allPods
}
