package map_utils

// MergeMaps merges two maps. If duplicate keys exist, the second map's value overwrites the first.
func MergeMaps[K comparable, V any](map1, map2 map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range map1 {
		result[k] = v
	}
	for k, v := range map2 {
		result[k] = v
	}
	return result
}

// InvertMap inverts a map, swapping keys and values. Assumes unique values.
func InvertMap[K comparable, V comparable](m map[K]V) map[V]K {
	inverted := make(map[V]K)
	for k, v := range m {
		inverted[v] = k
	}
	return inverted
}