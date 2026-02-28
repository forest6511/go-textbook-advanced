package main

import (
	"fmt"
	"time"
)

// Functional Options pattern
type ServerConfig struct {
	host    string
	port    int
	timeout time.Duration
}

type Option func(*ServerConfig)

func WithHost(host string) Option {
	return func(c *ServerConfig) {
		c.host = host
	}
}

func WithPort(port int) Option {
	return func(c *ServerConfig) {
		c.port = port
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *ServerConfig) {
		c.timeout = d
	}
}

func NewServer(opts ...Option) *ServerConfig {
	cfg := &ServerConfig{
		host:    "localhost",
		port:    8080,
		timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// Strategy pattern
type SortStrategy func([]int)

func bubbleSort(data []int) {
	n := len(data)
	for i := range n {
		for j := 0; j < n-i-1; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

func Sorter(data []int, strategy SortStrategy) {
	strategy(data)
}

// Repository pattern
type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	FindByID(id int) (*User, error)
	Save(user *User) error
}

type InMemoryUserRepository struct {
	users  map[int]*User
	nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[int]*User)}
}

func (r *InMemoryUserRepository) FindByID(id int) (*User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (r *InMemoryUserRepository) Save(user *User) error {
	r.nextID++
	user.ID = r.nextID
	r.users[user.ID] = user
	return nil
}

func main() {
	srv := NewServer(
		WithHost("0.0.0.0"),
		WithPort(9090),
		WithTimeout(60*time.Second),
	)
	fmt.Printf("server: %s:%d timeout=%s\n", srv.host, srv.port, srv.timeout)

	data := []int{64, 34, 25, 12, 22, 11, 90}
	Sorter(data, bubbleSort)
	fmt.Println("sorted:", data)

	repo := NewInMemoryUserRepository()
	user := &User{Name: "Alice"}
	_ = repo.Save(user)
	found, _ := repo.FindByID(user.ID)
	fmt.Println("found user:", found.Name)
}
