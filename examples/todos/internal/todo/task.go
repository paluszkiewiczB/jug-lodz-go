package todo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ID string

func (i ID) String() string {
	return "ID[" + string(i) + "]"
}

var (
	seed = rand.Int63()
	src  = rand.NewSource(seed)
)

func RandomID() ID {
	// standard library does not implement UUIDs :(
	return ID(strconv.Itoa(int(src.Int63())))
}

type Task struct {
	ID       ID
	Title    string
	Deadline *time.Time
	Done     bool
}

type Storage interface {
	Upsert(ctx context.Context, t Task) (stored Task, err error)
	List(ctx context.Context, f *TaskFilter) ([]Task, error)
}

type TaskFilter struct {
	ID *ID
}

type Handler struct {
	s Storage
}

func NewHandler(s Storage) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Create(ctx context.Context, cmd CreateTask) (Task, error) {
	var t = Task{
		ID:       RandomID(),
		Title:    cmd.Title,
		Deadline: cmd.Deadline,
		Done:     false,
	}

	stored, err := h.s.Upsert(ctx, t)
	if err != nil {
		return Task{}, fmt.Errorf("upserting the task: %v, %w", t, err)
	}

	return stored, nil
}

func (h *Handler) Toggle(ctx context.Context, id ID) (Task, error) {
	tasks, err := h.s.List(ctx, &TaskFilter{ID: &id})
	if err != nil {
		return Task{}, fmt.Errorf("listing task to toggle by id: %s, %w", id, err)
	}

	if len(tasks) == 0 {
		return Task{}, fmt.Errorf("by id: %s, %w", id, ErrTaskNotFound)
	}

	if v := len(tasks); v > 1 {
		return Task{}, fmt.Errorf("expected to find one task by id: %s, found: %d", id, v)
	}

	found := tasks[0]
	found.Done = !found.Done

	stored, err := h.s.Upsert(ctx, found)
	if err != nil {
		return Task{}, fmt.Errorf("upserting task: %s after toggling it: ")
	}

	return stored, nil
}

type CreateTask struct {
	Title string
	// Deadline is optional - nil means there is no deadline
	Deadline *time.Time
}

var ErrTaskNotFound = errors.New("task not found")
