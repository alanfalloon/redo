package tgt

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockT struct {
	mock.Mock
}
func (m *MockT) Watch(p string, r Observer) {
	m.Mock.Called(p, r)
	r <- expected_update
}

type MockFactory struct {
	mock.Mock
}
func (m *MockFactory) Create(p string) T {
	args := m.Mock.Called(p)
	return args.Get(0).(*MockT)
}

var expected_update = Update{"foo/../bar", "bar", UPDATED}

func TestWatch(t *testing.T) {
	assert := assert.New(t)
	tobs := make(chan Update)
	tfac := new(MockFactory)
	ttgt := new(MockT)
	tfac.On("Create", "foo/../bar").Return(ttgt)
	ttgt.On("Watch", "foo/../bar", tobs).Return()
	Watch("foo/../bar", tobs, tfac)
	assert.Equal(<-tobs, expected_update)
	tfac.Mock.AssertExpectations(t)
	ttgt.Mock.AssertExpectations(t)
}


func TestClaim(t *testing.T) {
	assert := assert.New(t)
	tobs := make(chan Update)
	tfac := new(MockFactory)
	ttgt := new(MockT)
	ttgt.On("Watch", expected_update.Alias, tobs).Return()
	Claim(expected_update.Alias, ttgt)
	Watch("foo/../bar", tobs, tfac)
	assert.Equal(<-tobs, expected_update)
	tfac.Mock.AssertExpectations(t)
	ttgt.Mock.AssertExpectations(t)
}
