// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	djikstra "github.com/dnnrly/layli/pathfinder/dijkstra"

	mock "github.com/stretchr/testify/mock"
)

// PathFinder is an autogenerated mock type for the PathFinder type
type PathFinder struct {
	mock.Mock
}

// AddConnection provides a mock function with given fields: from, cost, to
func (_m *PathFinder) AddConnection(from djikstra.Point, cost djikstra.CostFunction, to ...djikstra.Point) {
	_va := make([]interface{}, len(to))
	for _i := range to {
		_va[_i] = to[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, from, cost)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// BestPath provides a mock function with given fields:
func (_m *PathFinder) BestPath() ([]djikstra.Point, error) {
	ret := _m.Called()

	var r0 []djikstra.Point
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]djikstra.Point, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []djikstra.Point); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]djikstra.Point)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPathFinder creates a new instance of PathFinder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPathFinder(t interface {
	mock.TestingT
	Cleanup(func())
}) *PathFinder {
	mock := &PathFinder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}