package registry

import "github.com/mfine30/prodda/domain"

type ProdRegistry interface {
	All() ([]*domain.Prod, error)
	Add(p *domain.Prod) error
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
	r.prods = append(r.prods, p)
	return nil
}
