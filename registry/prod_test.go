package registry_test

import (
	"github.com/mfine30/prodda/domain"
	"github.com/mfine30/prodda/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InMemoryProdRegistry", func() {
	It("is empty on initialization", func() {
		r := registry.NewInMemoryProdRegistry()
		allProds, err := r.All()
		Expect(err).NotTo(HaveOccurred())
		Expect(allProds).NotTo(BeNil())
		Expect(len(allProds)).To(Equal(0))
	})

	It("stores a prod when added", func() {
		prod := &domain.Prod{}
		r := registry.NewInMemoryProdRegistry()
		r.Add(prod)
		allProds, err := r.All()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(allProds)).To(Equal(1))
		Expect(allProds[0]).To(Equal(prod))
	})

	It("assigns a new ID to prod when added", func() {
		prod := &domain.Prod{}
		r := registry.NewInMemoryProdRegistry()
		r.Add(prod)
		Expect(prod.ID).ToNot(Equal(0))
	})
})
