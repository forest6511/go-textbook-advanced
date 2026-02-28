package main

import (
	"fmt"
	"time"
)

// Age is a named type for age validation
type Age int

func (a Age) Validate() error {
	if a < 0 || a > 150 {
		return fmt.Errorf("age must be between 0 and 150, got %d", a)
	}
	return nil
}

// RewardPoints is a named type for reward points
type RewardPoints int

// Person demonstrates new(expr) from Go 1.26
type Person struct {
	Name string `json:"name"`
	Age  *int   `json:"age,omitempty"`
}

func yearsSince(t time.Time) int {
	return int(time.Since(t).Hours() / (365.25 * 24))
}

func newPerson(name string, born time.Time) Person {
	return Person{
		Name: name,
		Age:  new(yearsSince(born)), // Go 1.26: new(expr)
	}
}

// Config demonstrates new(expr) for optional fields
type Config struct {
	Timeout     *int
	MaxRetries  *int
	ServiceName *string
}

func defaultConfig() Config {
	return Config{
		Timeout:     new(30),
		MaxRetries:  new(3),
		ServiceName: new("payment-service"),
	}
}

func main() {
	age := Age(25)
	if err := age.Validate(); err != nil {
		fmt.Println("invalid:", err)
		return
	}
	fmt.Println("age valid:", age)

	born := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	p := newPerson("Alice", born)
	fmt.Printf("person: %s, age: %d\n", p.Name, *p.Age)

	cfg := defaultConfig()
	fmt.Printf("timeout: %d, retries: %d, service: %s\n",
		*cfg.Timeout, *cfg.MaxRetries, *cfg.ServiceName)
}
