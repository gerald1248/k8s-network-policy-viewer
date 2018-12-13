package main

import (
	"os"
	"strings"
)

const blacklistEnv = "NETWORK_POLICY_VIEWER_BLACKLIST"

func onBlacklist(namespace string) bool {
	blacklist := os.Getenv(blacklistEnv)
	return onBlacklistNoEnv(namespace, blacklist)
}

func onBlacklistNoEnv(namespace string, blacklist string) bool {
	if strings.TrimSpace(namespace) == "" || strings.TrimSpace(blacklist) == "" {
		return false
	}

	items := strings.Split(blacklist, ",")

	for _, item := range items {
		if strings.Contains(namespace, item) {
			return true
		}
	}
	return false
}
