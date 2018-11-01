package main

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
