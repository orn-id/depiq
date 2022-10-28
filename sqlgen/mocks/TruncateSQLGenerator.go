// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import exp "github.com/orn-id/depiq/v9/exp"
import mock "github.com/stretchr/testify/mock"
import sb "github.com/orn-id/depiq/v9/internal/sb"

// TruncateSQLGenerator is an autogenerated mock type for the TruncateSQLGenerator type
type TruncateSQLGenerator struct {
	mock.Mock
}

// Dialect provides a mock function with given fields:
func (_m *TruncateSQLGenerator) Dialect() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Generate provides a mock function with given fields: b, clauses
func (_m *TruncateSQLGenerator) Generate(b sb.SQLBuilder, clauses exp.TruncateClauses) {
	_m.Called(b, clauses)
}
