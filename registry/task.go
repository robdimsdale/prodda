package registry

import (
	"math/rand"

	"github.com/mfine30/prodda/domain"
)

type TaskRegistry interface {

	// All returns all the tasks known to the registry.
	// If no tasks exist, but the interrogation was otherwise successful,
	// the returned error will be nil.
	All() ([]domain.Task, error)

	// Add adds the task to current tasks known to the registry.
	// Add is responsible for assiging a unique ID to the provided tasks.
	Add(p domain.Task) error

	// ByID will return error if there is an error retriving a task which exists.
	// If the execution completes without error, and task is not found,
	// both the returned error and task will be nil.
	ByID(ID uint) (domain.Task, error)

	// Update will return an error if the task does not exist.
	// Callers are expected to first verify that the task exists,
	// e.g. via ByID.
	Update(task domain.Task) (domain.Task, error)

	// Remove will return an error if the task does not exist.
	// Callers are expected to first determine if the task exists,
	// e.g. via ByID.
	Remove(task domain.Task) error
}

type InMemoryTaskRegistry struct {
	tasks []domain.Task
}

func NewInMemoryTaskRegistry() TaskRegistry {
	return &InMemoryTaskRegistry{
		tasks: []domain.Task{},
	}
}

func (r InMemoryTaskRegistry) All() ([]domain.Task, error) {
	return r.tasks, nil
}

func (r *InMemoryTaskRegistry) Add(p domain.Task) error {
	err := p.SetID(r.uniqueRandomID())
	if err != nil {
		return err
	}

	r.tasks = append(r.tasks, p)
	return nil
}

func (r InMemoryTaskRegistry) uniqueRandomID() uint {
	newID := uint(rand.Uint32())
	existingTask, _ := r.ByID(newID)
	for existingTask != nil {
		newID := uint(rand.Uint32())
		existingTask, _ = r.ByID(newID)
	}
	return newID
}

func (r InMemoryTaskRegistry) ByID(ID uint) (domain.Task, error) {
	_, found, err := r.byID(ID)
	if err != nil {
		return nil, err
	}
	return found, nil
}

func (r InMemoryTaskRegistry) byID(ID uint) (int, domain.Task, error) {
	allTasks, err := r.All()
	if err != nil {
		return 0, nil, err
	}

	for i, p := range allTasks {
		if p.ID() == ID {
			return i, p, nil
		}
	}
	return 0, nil, nil
}

func (r *InMemoryTaskRegistry) Update(task domain.Task) (domain.Task, error) {
	found, err := r.ByID(task.ID())
	if err != nil {
		return nil, err
	}

	found.SetSchedule(task.Schedule())

	return found, nil
}

func (r *InMemoryTaskRegistry) Remove(task domain.Task) error {
	i, _, err := r.byID(task.ID())
	if err != nil {
		return err
	}

	r.tasks[i] = nil // explicitly set to nil to avoid memory leaks
	r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)

	return nil
}
