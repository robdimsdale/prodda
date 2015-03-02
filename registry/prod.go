package registry

import (
	"math/rand"

	"github.com/mfine30/prodda/domain"
)

type ProdRegistry interface {

	// All returns all the prods known to the registry.
	// If no prods exist, but the interrogation was otherwise successful,
	// the returned error will be nil.
	All() ([]*domain.Prod, error)

	// Add adds the prod to current prods in memory.
	// Add is responsible for assiging a unique ID to the provided prod.
	Add(p *domain.Prod) error

	// ByID will return error if there is an error retriving a prod which exists.
	// If the execution completes without error, and prod is not found,
	// both the returned error and prod will be nil.
	ByID(ID int) (*domain.Prod, error)

	// Update will return an error if the prod does not exist.
	// Callers are expected to first verify that the prod exists,
	// e.g. via ByID.
	Update(prod *domain.Prod) (*domain.Prod, error)

	// Remove will return an error if the prod does not exist.
	// Callers are expected to first determine if the prod exists,
	// e.g. via ByID.
	Remove(prod *domain.Prod) error
}

type InMemoryProdRegistry struct {
	prods []*domain.Prod
}

func NewInMemoryProdRegistry() ProdRegistry {
	return &InMemoryProdRegistry{
		prods: []*domain.Prod{},
	}
}

func (r InMemoryProdRegistry) All() ([]*domain.Prod, error) {
	return r.prods, nil
}

func (r *InMemoryProdRegistry) Add(p *domain.Prod) error {

	newID := rand.Int()
	existingProd, _ := r.ByID(newID)
	for existingProd != nil {
		newID := rand.Int()
		existingProd, _ = r.ByID(newID)
	}

	p.ID = newID
	r.prods = append(r.prods, p)
	return nil
}

func (r InMemoryProdRegistry) ByID(ID int) (*domain.Prod, error) {
	_, found, err := r.byID(ID)
	if err != nil {
		return nil, err
	}
	return found, nil
}

func (r InMemoryProdRegistry) byID(ID int) (int, *domain.Prod, error) {
	allProds, err := r.All()
	if err != nil {
		return 0, nil, err
	}

	for i, p := range allProds {
		if p.ID == ID {
			return i, p, nil
		}
	}
	return 0, nil, nil
}

func (r *InMemoryProdRegistry) Update(prod *domain.Prod) (*domain.Prod, error) {
	found, err := r.ByID(prod.ID)
	if err != nil {
		return nil, err
	}

	found.Schedule = prod.Schedule

	return found, nil
}

func (r *InMemoryProdRegistry) Remove(prod *domain.Prod) error {
	i, _, err := r.byID(prod.ID)
	if err != nil {
		return err
	}

	r.prods[i] = nil // explicitly set to nil to avoid memory leaks
	r.prods = append(r.prods[:i], r.prods[i+1:]...)

	return nil
}
