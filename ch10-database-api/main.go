package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

// User domain model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRepository interface (Repository pattern)
type UserRepository interface {
	FindByID(id int) (*User, error)
	Save(user *User) error
}

// InMemoryUserRepository is an in-memory implementation
type InMemoryUserRepository struct {
	users  map[int]*User
	nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[int]*User)}
}

var ErrNotFound = errors.New("not found")

func (r *InMemoryUserRepository) FindByID(id int) (*User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %d: %w", id, ErrNotFound)
	}
	return u, nil
}

func (r *InMemoryUserRepository) Save(user *User) error {
	r.nextID++
	user.ID = r.nextID
	r.users[user.ID] = user
	return nil
}

// Handlers
var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Error("json encode failed", "err", err)
	}
}

func getUserHandler(repo UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		user, err := repo.FindByID(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}

func createUserHandler(repo UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// io.ReadAll is ~2x faster in Go 1.26
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "read failed"})
			return
		}
		var user User
		if err := json.Unmarshal(body, &user); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		if err := repo.Save(&user); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "save failed"})
			return
		}
		writeJSON(w, http.StatusCreated, user)
	}
}

func newServer(repo UserRepository) http.Handler {
	// CrossOriginProtection for CSRF (Go 1.25)
	protection := http.NewCrossOriginProtection()

	mux := http.NewServeMux()
	// Go 1.22 enhanced routing
	mux.HandleFunc("GET /api/users/{id}", getUserHandler(repo))
	mux.HandleFunc("POST /api/users", createUserHandler(repo))

	return protection.Handler(mux)
}

func main() {
	repo := NewInMemoryUserRepository()
	handler := newServer(repo)
	logger.Info("server starting", slog.String("addr", ":8080"))
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("server error:", err)
	}
}
