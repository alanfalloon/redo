package bottomup

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/alanfalloon/redo/bottomup/tgt"
)

type MockTgtsI struct {
	mock.Mock
	tgt tgt.T
}
func (m *MockTgtsI) Watch(path string, r tgt.Observer, fac tgt.Factory) {
	args := m.Mock.Called(path, r, fac)
	if args.Bool(0) {
		m.tgt = fac.Create(path)
		m.tgt.Watch(path, r)
	}
}
func (m *MockTgtsI) Claim(path string, tgt tgt.T) {
	m.Mock.Called(path, tgt)
}

type MockT struct {
	mock.Mock
}
func (m *MockT) Watch(p string, r tgt.Observer) {
	args := m.Mock.Called(p, r)
	go func(u tgt.Update){r<-u}(args.Get(0).(tgt.Update))
}

type MockFactory struct {
	mock.Mock
}
func (m *MockFactory) Create(p string) tgt.T {
	args := m.Mock.Called(p)
	return args.Get(0).(*MockT)
}

func TestAlias(t *testing.T) {
	mtgts := new(MockTgtsI)
	mfac := new(MockFactory)
	mtgt := new(MockT)
	aliasfac := mkaliasfac(mtgts, mfac)
	obs := make(chan tgt.Update)

	mtgts.Mock.On("Watch", "foo/bar", obs, aliasfac).Return(true).Once()
	mfac.Mock.On("Create", "foo/bar").Return(mtgt).Once()
	mtgts.Mock.On("Claim", "foo/bar", mtgt).Return().Once()
	exp_u := tgt.Update{"foo/bar", "foo/bar", tgt.MISSING}
	mtgt.Mock.On("Watch", "foo/bar", obs).Return(exp_u).Once()

	mtgts.Watch("foo/bar", obs, aliasfac)

	u := <- obs
	assert.Equal(t, exp_u, u)
	
	mock.AssertExpectationsForObjects(t, mtgts.Mock, mfac.Mock, mtgt.Mock)
}
