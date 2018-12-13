package main

import (
	"testing"
)

func TestBlacklistOn(t *testing.T) {
	namespace := "bob"
	blacklist := "alice,bob,eve"
	if !onBlacklistNoEnv(namespace, blacklist) {
		t.Errorf("Lookup '%s' in '%v' must return true", namespace, blacklist)
	}
}

func TestBlacklistNotOn(t *testing.T) {
	namespace := "bob"
	blacklist := "alice,eve"
	if onBlacklistNoEnv(namespace, blacklist) {
		t.Errorf("Lookup '%s' in '%v' must return false", namespace, blacklist)
	}
}
