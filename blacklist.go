package main

import (
	"os"
	"strings"
)

const blacklistEnv = "NETWORK_POLICY_VIEWER_BLACKLIST"

func onBlacklist(namespace string) bool {
	if len(namespace) == 0 {
		return false
	}

	items := strings.Split(os.Getenv(blacklistEnv), ",")
	for _, item := range items {
		if strings.Contains(namespace, item) {
			return true
		}
	}
	return false
}
