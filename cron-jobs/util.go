package main

import (
	"fmt"
	"os"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	fmt.Printf("Warning: %s not set, using fallback value: %s\n", key, fallback)
	return fallback
}


func toSet(list []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range list {
		set[v] = true
	}
	return set
}
