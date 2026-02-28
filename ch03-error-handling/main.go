package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

// QueryError is a custom error type
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query %q: %v", e.Query, e.Err)
}

func (e *QueryError) Unwrap() error {
	return e.Err
}

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func runQuery(query string) error {
	if query == "" {
		return &QueryError{
			Query: query,
			Err:   errors.New("empty query"),
		}
	}
	return nil
}

func handleError(err error) {
	// Recommended (Go 1.26): errors.AsType[T]() - type-safe, no pre-declaration needed
	if qErr, ok := errors.AsType[*QueryError](err); ok {
		logger.Error("query failed (AsType)", "query", qErr.Query, "err", qErr.Err)
	}

	// Compatibility (Go 1.25 and earlier): errors.As
	var qErr2 *QueryError
	if errors.As(err, &qErr2) {
		logger.Error("query failed (As)", "query", qErr2.Query, "err", qErr2.Err)
	}
}

// Sentinel errors
var (
	ErrNotFound = errors.New("not found")
	ErrTimeout  = errors.New("timeout")
)

// errors.Join example (Go 1.20+)
func validate(name, email string) error {
	var errs []error
	if name == "" {
		errs = append(errs, errors.New("name is required"))
	}
	if email == "" {
		errs = append(errs, errors.New("email is required"))
	}
	return errors.Join(errs...)
}

func main() {
	err := runQuery("")
	if err != nil {
		handleError(err)
	}

	err2 := validate("", "")
	if err2 != nil {
		fmt.Println("validation error:", err2)
	}

	// errors.Is with sentinel
	wrappedNotFound := fmt.Errorf("user lookup: %w", ErrNotFound)
	fmt.Println("is ErrNotFound:", errors.Is(wrappedNotFound, ErrNotFound))
}
