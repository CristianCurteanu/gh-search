package slices

func MapSlice[I any, O any](input []I, cb func(I) O) []O {
	var output []O = make([]O, 0, len(input))
	for _, i := range input {
		output = append(output, cb(i))
	}

	return output
}
