package bottomup

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTgtsI struct {
	mock.Mock
}
func (m *MockTgtsI) Watch(path string, resp tgt.Observer, fac Factory) {
	m.Mock.Called(path, resp, fac)
}
func (m *MockTgtsI) Claim(path string, tgt tgt.T) {
	m.Mock.Called(path, tgt)
}

type MockT struct {
	mock.Mock
}
func (m *MockT) Watch(p string, r tgt.Observer) {
	m.Mock.Called(p, r)
}

type MockFactory struct {
	mock.Mock
}
func (m *MockFactory) Create(p string) tgt.T {
	args := m.Mock.Called(p)
	return args.Get(0).(*MockT)
}

func TestAlias(t *testing.T) {
	assert := assert.New(t)
}
