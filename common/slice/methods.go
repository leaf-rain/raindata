package slice

func ContainsElementString(slice []string, element string) int {
	for index, value := range slice {
		if value == element {
			return index
		}
	}
	return -1
}
