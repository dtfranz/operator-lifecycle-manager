// Code generated by counterfeiter. DO NOT EDIT.
package listers

import (
	"sync"

	v1a "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/rbac/v1"
)

type FakeClusterRoleBindingLister struct {
	GetStub        func(string) (*v1a.ClusterRoleBinding, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 string
	}
	getReturns struct {
		result1 *v1a.ClusterRoleBinding
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 *v1a.ClusterRoleBinding
		result2 error
	}
	ListStub        func(labels.Selector) ([]*v1a.ClusterRoleBinding, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		arg1 labels.Selector
	}
	listReturns struct {
		result1 []*v1a.ClusterRoleBinding
		result2 error
	}
	listReturnsOnCall map[int]struct {
		result1 []*v1a.ClusterRoleBinding
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClusterRoleBindingLister) Get(arg1 string) (*v1a.ClusterRoleBinding, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetStub
	fakeReturns := fake.getReturns
	fake.recordInvocation("Get", []interface{}{arg1})
	fake.getMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClusterRoleBindingLister) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeClusterRoleBindingLister) GetCalls(stub func(string) (*v1a.ClusterRoleBinding, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *FakeClusterRoleBindingLister) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClusterRoleBindingLister) GetReturns(result1 *v1a.ClusterRoleBinding, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *v1a.ClusterRoleBinding
		result2 error
	}{result1, result2}
}

func (fake *FakeClusterRoleBindingLister) GetReturnsOnCall(i int, result1 *v1a.ClusterRoleBinding, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 *v1a.ClusterRoleBinding
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 *v1a.ClusterRoleBinding
		result2 error
	}{result1, result2}
}

func (fake *FakeClusterRoleBindingLister) List(arg1 labels.Selector) ([]*v1a.ClusterRoleBinding, error) {
	fake.listMutex.Lock()
	ret, specificReturn := fake.listReturnsOnCall[len(fake.listArgsForCall)]
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		arg1 labels.Selector
	}{arg1})
	stub := fake.ListStub
	fakeReturns := fake.listReturns
	fake.recordInvocation("List", []interface{}{arg1})
	fake.listMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClusterRoleBindingLister) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeClusterRoleBindingLister) ListCalls(stub func(labels.Selector) ([]*v1a.ClusterRoleBinding, error)) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = stub
}

func (fake *FakeClusterRoleBindingLister) ListArgsForCall(i int) labels.Selector {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	argsForCall := fake.listArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClusterRoleBindingLister) ListReturns(result1 []*v1a.ClusterRoleBinding, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 []*v1a.ClusterRoleBinding
		result2 error
	}{result1, result2}
}

func (fake *FakeClusterRoleBindingLister) ListReturnsOnCall(i int, result1 []*v1a.ClusterRoleBinding, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	if fake.listReturnsOnCall == nil {
		fake.listReturnsOnCall = make(map[int]struct {
			result1 []*v1a.ClusterRoleBinding
			result2 error
		})
	}
	fake.listReturnsOnCall[i] = struct {
		result1 []*v1a.ClusterRoleBinding
		result2 error
	}{result1, result2}
}

func (fake *FakeClusterRoleBindingLister) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClusterRoleBindingLister) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ v1.ClusterRoleBindingLister = new(FakeClusterRoleBindingLister)
