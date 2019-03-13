package main

func selectNamespaces(selector *map[string]string, namespacePodMap *map[string][]string, namespaceLabelMap *map[string]map[string]string) []string {
	var namespaces []string
	for key := range *namespacePodMap {
		namespaces = append(namespaces, key)
	}

	// special case: empty map selects all pods
	if len(*selector) == 0 {
		return namespaces
	}

	var selectedNamespaces []string
	for _, namespace := range namespaces {
		labels := (*namespaceLabelMap)[namespace]
		for k, v := range *selector {
			if len(k) > 0 && labels[k] == v {
				selectedNamespaces = append(selectedNamespaces, namespace)
			}
		}
	}

	return selectedNamespaces
}
