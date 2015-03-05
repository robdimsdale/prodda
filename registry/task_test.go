package registry_test

import (
	"github.com/prodda/prodda/domain"
	"github.com/prodda/prodda/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InMemoryTaskRegistry", func() {
	It("is empty on initialization", func() {
		r := registry.NewInMemoryTaskRegistry()
		allTasks, err := r.All()
		Expect(err).NotTo(HaveOccurred())
		Expect(allTasks).NotTo(BeNil())
		Expect(len(allTasks)).To(Equal(0))
	})

	It("stores a task when added", func() {
		task := &domain.NoOpTask{}
		r := registry.NewInMemoryTaskRegistry()
		r.Add(task)
		allTasks, err := r.All()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(allTasks)).To(Equal(1))
		Expect(allTasks[0]).To(Equal(task))
	})

	It("assigns a new ID to task when added", func() {
		task := &domain.NoOpTask{}
		r := registry.NewInMemoryTaskRegistry()
		r.Add(task)
		Expect(task.ID).ToNot(Equal(0))
	})
})
