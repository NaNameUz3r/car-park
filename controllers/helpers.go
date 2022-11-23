package controllers

func contains(slice []int64, id int64) bool {
	for _, item := range slice {
		if item == id {
			return true
		}
	}
	return false
}
