package db

func getDefaultPageLine(page, line int) (int, int) {
	if page <= 0 {
		page = 1
	}

	if line <= 0 {
		line = 15
	}

	return page, line
}
