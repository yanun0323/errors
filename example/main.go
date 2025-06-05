package main

import (
	"context"
	builtinErrors "errors"
	"fmt"

	"github.com/yanun0323/errors"
)

var (
	ErrUserNotFound2 = builtinErrors.New("user not found")
	ErrUserNotFound  = errors.New("user not found")
)

func main() {
	ctx := context.Background()
	err := handleRequest(ctx, 0)

	if err != nil {
		fmt.Printf("\n=== short error message:\n%s\n", err)

		fmt.Printf("\n=== text:\n%v\n", err)

		fmt.Printf("\n===json:\n%#v\n", err)

		fmt.Printf("\n===colorized:\n%+v\n", err)
	}

	println(errors.Is(err, ErrUserNotFound))
}

func handleRequest(ctx context.Context, userID int) error {
	if err := processUser(ctx, userID); err != nil {
		return errors.Wrap(err, "process user").
			With("host", "db.example.com").
			With("port", 5432).
			With("timeout", "30s").
			With("func", "handleRequest")
	}

	return nil
}

// processUser handles user-related operations
func processUser(ctx context.Context, userID int) error {
	if err := validateUser(userID); err != nil {
		return errors.Errorf("user validation failed, err: %w", err).
			With("func", "processUser")
	}

	return nil
}

// validateUser validates if user exists
func validateUser(userID int) error {
	if userID == 0 {
		return errors.Errorf("root: %w", ErrUserNotFound).
			With("user_id", userID).
			With("table", "users").
			With("func", "validateUser")
	}

	return nil
}
