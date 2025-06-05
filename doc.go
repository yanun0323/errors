// Package errors provides a lightweight error handling library with support for
// error wrapping, formatting, and error chain checking.
//
// # DO NOT use fmt.Errorf to create errors, use errors.New or errors.Errorf instead.
//
// fmt.Errorf is not compatible with errors.Is and errors.As.
//
//	fmt.Errorf("failed to process user %d: %w", userID, err)
//
// is not the same as
//
//	errors.Errorf("failed to process user %d: %w", userID, err)
//
// # Example:
//
//	// Create an error
//	err := errors.New("user validation failed").
//		With("user_id", 12345).
//		With("email", "user@example.com").
//		With("attempt", 3)
//
//	// Wrap the error
//	err = errors.Wrap(err, "using wrap to wrap the error")
//	err = errors.Errorf("%w can also wrap the error", err)
//
//	// Format the error
//	message := err.Error()
//	textWithStack := errors.Format(err)
//	jsonWithStack := errors.FormatJSON(err)
//	colorizedWithStack := errors.FormatColorized(err)
//
//	// using Sprintf to format the error
//	message := fmt.Sprintf("%s", err)
//	textWithStack := fmt.Sprintf("%v", err)  // equal to errors.Format(err)
//	jsonWithStack := fmt.Sprintf("%#v", err)  // equal to errors.FormatJSON(err)
//	colorizedWithStack := fmt.Sprintf("%+v", err)  // equal to errors.FormatColorized(err)
//
//	// Check if the error is a specific error
//	if errors.Is(err, errors.New("user validation failed")) {
//		// handle the error
//	}
//
//	var validationErr AnErrorType
//	if errors.As(err, &validationErr) {
//		// handle the error
//	}
package errors
