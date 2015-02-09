// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/mfine30/prodda/timer"
)

type FakeTask struct {
	RunStub        func() error
	runMutex       sync.RWMutex
	runArgsForCall []struct{}
	runReturns struct {
		result1 error
	}
}

func (fake *FakeTask) Run() error {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct{}{})
	fake.runMutex.Unlock()
	if fake.RunStub != nil {
		return fake.RunStub()
	} else {
		return fake.runReturns.result1
	}
}

func (fake *FakeTask) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *FakeTask) RunReturns(result1 error) {
	fake.RunStub = nil
	fake.runReturns = struct {
		result1 error
	}{result1}
}

var _ prodda.Task = new(FakeTask)
