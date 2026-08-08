package store

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("todo not found")

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TodoStore struct {
	mu     sync.RWMutex
	items  map[int]Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		items:  make(map[int]Todo),
		nextID: 1,
	}
}

func (s *TodoStore) List() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]Todo, 0, len(s.items))
	for _, t := range s.items {
		todos = append(todos, t)
	}
	return todos
}

func (s *TodoStore) Get(id int) (Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	return t, nil
}

func (s *TodoStore) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := Todo{ID: s.nextID, Title: title, Done: false}
	s.items[t.ID] = t
	s.nextID++
	return t
}
