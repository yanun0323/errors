package errors

const (
	_emptyJSONString = "{}"
	_emptyString     = ""
)

// Format formats the error as a string
func Format(err error) string {
	if err == nil {
		return _emptyString
	}

	if err, ok := err.(*errorStack); ok {
		return err.formatText()
	}

	return err.Error()
}

// FormatJson formats the error as a JSON string
func FormatJson(err error) string {
	if err == nil {
		return _emptyJSONString
	}

	if err, ok := err.(*errorStack); ok {
		return err.formatJson()
	}

	return err.Error()
}

// FormatColorized formats the error as a colorized string
func FormatColorized(err error) string {
	if err == nil {
		return _emptyString
	}

	if err, ok := err.(*errorStack); ok {
		return err.formatColorized()
	}

	return err.Error()
}
