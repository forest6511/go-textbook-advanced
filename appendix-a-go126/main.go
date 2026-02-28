package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

// Demonstrate Go 1.26 new(expr)
type Config struct {
	Timeout     *int
	MaxRetries  *int
	ServiceName *string
}

func defaultConfig() Config {
	return Config{
		Timeout:     new(30),
		MaxRetries:  new(3),
		ServiceName: new("myservice"),
	}
}

// Demonstrate Go 1.26 self-referential generics
type Adder[A Adder[A]] interface {
	Add(A) A
}

type Vec2 struct{ X, Y float64 }

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

func algo[A Adder[A]](x, y A) A {
	return x.Add(y)
}

// Demonstrate Go 1.26 errors.AsType[T]()
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code %d: %s", e.Code, e.Message)
}

// Demonstrate Go 1.25 WaitGroup.Go()
func parallelWork(tasks []int) []int {
	results := make([]int, len(tasks))
	var wg sync.WaitGroup
	for i, t := range tasks {
		i, t := i, t
		wg.Go(func() {
			results[i] = t * t
		})
	}
	wg.Wait()
	return results
}

// Demonstrate Go 1.26 slog.NewMultiHandler
func newMultiLogger() *slog.Logger {
	console := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	file := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(slog.NewMultiHandler(console, file))
}

func main() {
	// new(expr)
	cfg := defaultConfig()
	fmt.Printf("config: timeout=%d retries=%d service=%s\n",
		*cfg.Timeout, *cfg.MaxRetries, *cfg.ServiceName)

	// Self-referential generics
	a := Vec2{1, 2}
	b := Vec2{3, 4}
	result := algo(a, b)
	fmt.Printf("Vec2 add: {%.1f, %.1f}\n", result.X, result.Y)

	// errors.AsType[T]()
	err := fmt.Errorf("wrapped: %w", &AppError{Code: 404, Message: "not found"})
	if appErr, ok := errors.AsType[*AppError](err); ok {
		fmt.Printf("AppError code=%d msg=%s\n", appErr.Code, appErr.Message)
	}

	// WaitGroup.Go()
	tasks := []int{1, 2, 3, 4, 5}
	squares := parallelWork(tasks)
	fmt.Println("squares:", squares)

	// NewMultiHandler
	logger := newMultiLogger()
	logger.Info("Go 1.26 features demo", slog.Time("at", time.Now()))
}
