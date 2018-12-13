package main

import (
	"testing"
)

func TestSelectPodsMatchExpressionsIn(t *testing.T) {
	// signature: func selectPods(namespace string, selector *Selector, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string

	namespace := "default"

	// assemble selector

	// blank label selector
	selector := Selector{}
	selector.MatchLabels = map[string]string{}

	// match expressions
	values := []string{"alice", "bob", "eve"}
	matchExpression := LabelSelectorRequirement{}
	matchExpression.Key = "app"
	matchExpression.Operator = "In" // <- operator
	matchExpression.Values = values
	requirements := []*LabelSelectorRequirement{}
	requirements = append(requirements, &matchExpression)
	selector.MatchExpressions = requirements

	// assemble namespacePodMap
	namespacePodMap := map[string][]string{}
	nestedArray := []string{}
	nestedArray = append(nestedArray, "alice")
	nestedArray = append(nestedArray, "bob")
	namespacePodMap["default"] = nestedArray

	// assemble podLabelMap
	labelMapAlice := map[string]string{}
	labelMapAlice["app"] = "alice"
	labelMapBob := map[string]string{}
	labelMapBob["app"] = "bob"
	labelMapEve := map[string]string{}
	labelMapEve["app"] = "eve"
	podLabelMap := map[string]map[string]string{}

	podLabelMap["alice"] = labelMapAlice
	podLabelMap["bob"] = labelMapBob
	podLabelMap["eve"] = labelMapEve

	pods := selectPods(namespace, &selector, &namespacePodMap, &podLabelMap)

	if len(pods) != 2 {
		t.Errorf("Operator In: must find two matching pods - got array=%v", pods)
	}
}

func TestSelectPodsMatchExpressionsNotIn(t *testing.T) {
	// signature: func selectPods(namespace string, selector *Selector, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string

	namespace := "default"

	// assemble selector

	// blank label selector
	selector := Selector{}
	selector.MatchLabels = map[string]string{}

	// match expressions
	values := []string{"eve"}
	matchExpression := LabelSelectorRequirement{}
	matchExpression.Key = "app"
	matchExpression.Operator = "NotIn" // <- operator
	matchExpression.Values = values
	requirements := []*LabelSelectorRequirement{}
	requirements = append(requirements, &matchExpression)
	selector.MatchExpressions = requirements

	// assemble namespacePodMap
	namespacePodMap := map[string][]string{}
	nestedArray := []string{}
	nestedArray = append(nestedArray, "alice")
	nestedArray = append(nestedArray, "bob")
	namespacePodMap["default"] = nestedArray

	// assemble podLabelMap
	labelMapAlice := map[string]string{}
	labelMapAlice["app"] = "alice"
	labelMapBob := map[string]string{}
	labelMapBob["app"] = "bob"
	labelMapEve := map[string]string{}
	labelMapEve["app"] = "eve"
	podLabelMap := map[string]map[string]string{}

	podLabelMap["alice"] = labelMapAlice
	podLabelMap["bob"] = labelMapBob
	podLabelMap["eve"] = labelMapEve

	pods := selectPods(namespace, &selector, &namespacePodMap, &podLabelMap)

	if len(pods) != 2 {
		t.Errorf("OperatorNotIn: must find two matching pods - got array=%v", pods)
	}
}

func TestSelectPodsMatchExpressionsExists(t *testing.T) {
	// signature: func selectPods(namespace string, selector *Selector, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string

	namespace := "default"

	// assemble selector

	// blank label selector
	selector := Selector{}
	selector.MatchLabels = map[string]string{}

	// match expressions
	matchExpression := LabelSelectorRequirement{}
	matchExpression.Key = "app"
	matchExpression.Operator = "Exists" // <- operator
	requirements := []*LabelSelectorRequirement{}
	requirements = append(requirements, &matchExpression)
	selector.MatchExpressions = requirements

	// assemble namespacePodMap
	namespacePodMap := map[string][]string{}
	nestedArray := []string{}
	nestedArray = append(nestedArray, "alice")
	nestedArray = append(nestedArray, "bob")
	nestedArray = append(nestedArray, "eve")
	namespacePodMap["default"] = nestedArray

	// assemble podLabelMap
	labelMapAlice := map[string]string{}
	labelMapAlice["app"] = "alice"
	labelMapBob := map[string]string{}
	labelMapBob["app"] = "bob"
	labelMapEve := map[string]string{}
	labelMapEve["app"] = "eve"
	podLabelMap := map[string]map[string]string{}

	podLabelMap["alice"] = labelMapAlice
	podLabelMap["bob"] = labelMapBob
	podLabelMap["eve"] = labelMapEve

	pods := selectPods(namespace, &selector, &namespacePodMap, &podLabelMap)

	if len(pods) != 3 {
		t.Errorf("Exists: must find three matching pods - got array=%v", pods)
	}
}

func TestSelectPodsMatchExpressionsDoesNotExist(t *testing.T) {
	// signature: func selectPods(namespace string, selector *Selector, namespacePodMap *map[string][]string, podLabelMap *map[string]map[string]string) []string

	namespace := "default"

	// assemble selector

	// blank label selector
	selector := Selector{}
	selector.MatchLabels = map[string]string{}

	// match expressions
	matchExpression := LabelSelectorRequirement{}
	matchExpression.Key = "app"
	matchExpression.Operator = "DoesNotExist" // <- operator
	requirements := []*LabelSelectorRequirement{}
	requirements = append(requirements, &matchExpression)
	selector.MatchExpressions = requirements

	// assemble namespacePodMap
	namespacePodMap := map[string][]string{}
	nestedArray := []string{}
	nestedArray = append(nestedArray, "alice")
	nestedArray = append(nestedArray, "bob")
	nestedArray = append(nestedArray, "eve")
	namespacePodMap["default"] = nestedArray

	// assemble podLabelMap
	labelMapEve := map[string]string{}
	labelMapEve["app"] = "eve"
	podLabelMap := map[string]map[string]string{}
	podLabelMap["eve"] = labelMapEve

	pods := selectPods(namespace, &selector, &namespacePodMap, &podLabelMap)

	if len(pods) != 2 {
		t.Errorf("DoesNotExist: must find two matching pods - got array=%v", pods)
	}
}
