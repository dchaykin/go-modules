package datamodel

func SliceContains(list []string, value string) bool {
	for i := range list {
		if list[i] == value {
			return true
		}
	}
	return false
}
