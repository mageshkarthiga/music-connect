package main

func toSet(list []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range list {
		set[v] = true
	}
	return set
}
