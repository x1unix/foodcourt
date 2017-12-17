package orders

func FormatOrdersReport(statRowsPtr *[]UserOrderCounter) *map[int] []int {
	rows := *statRowsPtr
	output := make(map[int] []int)

	for _, val := range rows {
		if output[val.UserId] == nil {
			// Workaround for slices :(
			// If you have better solution - u r welcome
			output[val.UserId] = make([]int, 0)
		}

		output[val.UserId] = append(output[val.UserId], val.Date)
	}

	return &output
}