package registry

import (
	"fmt"
	"math/rand"

	"github.com/mfine30/prodda/domain"
)

type ProdRegistry interface {
	All() ([]*domain.Prod, error)

	// Add adds the prod to current prods in memory
	// Add is responsible for assiging a unique ID to the provided prod
	Add(p *domain.Prod) error
	ByID(ID int) (*domain.Prod, error)
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

// Add adds the prod to current prods in memory
// Add is responsible for assiging a unique ID to the provided prod
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
	allProds, err := r.All()
	if err != nil {
		return nil, err
	}

	for _, p := range allProds {
		if p.ID == ID {
			return p, nil
		}
	}
	return nil, fmt.Errorf("No prod found for id :%d", ID)
}
