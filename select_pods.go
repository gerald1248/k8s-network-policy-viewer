package main

func selectPods(namespace string, selector *Selector, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string {
	// special case: empty map selects all pods
	if selectorIsEmpty(selector) {
		//fmt.Printf("Selector empty: %v\n", *selector)
		return (*namespacePodMap)[namespace]
	} else {
		//fmt.Printf("Selector not empty: %v\n", *selector)
	}

	var selectedPods []string

	// select pods matching at least one label pair
	for _, pod := range (*namespacePodMap)[namespace] {
		labels := (*podLabelMap)[pod]
		for k, v := range (*selector).MatchLabels {
			if len(k) > 0 && labels[k] == v {
				selectedPods = append(selectedPods, pod)
			}
		}
	}

	// append pods matching any MatchExpressions
	for _, requirement := range (*selector).MatchExpressions {
		switch requirement.Operator {
		case "In":
			for _, pod := range (*namespacePodMap)[namespace] {
				for k, v := range (*podLabelMap)[pod] {
					if k == requirement.Key {
						for _, item := range requirement.Values {
							if v == item {
								selectedPods = append(selectedPods, pod)
								break
							}
						}
					}
				}
			}
		case "NotIn":
			for _, pod := range (*namespacePodMap)[namespace] {
				found := false
				for k, v := range (*podLabelMap)[pod] {
					if k == requirement.Key {
						for _, item := range requirement.Values {
							if v == item {
								found = true
							}
						}
					}
				}
				if found == false {
					selectedPods = append(selectedPods, pod)
				}
			}
		case "Exists":
			for _, pod := range (*namespacePodMap)[namespace] {
				for k, _ := range (*podLabelMap)[pod] {
					if k == requirement.Key {
						selectedPods = append(selectedPods, pod)
						continue
					}
				}
			}
		case "DoesNotExist":
			for _, pod := range (*namespacePodMap)[namespace] {
				found := false
				for k, _ := range (*podLabelMap)[pod] {
					if k == requirement.Key {
						found = true
					}
				}
				if found == false {
					selectedPods = append(selectedPods, pod)
				}
			}
		}
	}

	return selectedPods
}

func selectPodsAcrossNamespaces(namespaces *[]string, selector *Selector, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string {
	var allPods []string
	for _, namespace := range *namespaces {
		selectedPods := selectPods(namespace, selector, namespacePodMap, podLabelMap)
		allPods = append(allPods, selectedPods...)
	}
	return allPods
}

func selectorIsEmpty(selector *Selector) bool {
	matchLabelsEmpty := (*selector).MatchLabels == nil || len((*selector).MatchLabels) == 0
	matchExpressionsEmpty := (*selector).MatchExpressions == nil || len((*selector).MatchExpressions) == 0
	return matchLabelsEmpty && matchExpressionsEmpty
}
