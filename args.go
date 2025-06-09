package errors

func makeArgs(funcName string, args ...any) []attr {
	attrs := make([]attr, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			break
		}

		attrs = append(attrs, attr{
			Function: funcName,
			Key:      k,
			Value:    args[i+1],
		})
	}

	return attrs
}
