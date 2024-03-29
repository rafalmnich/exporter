// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"

import mock "github.com/stretchr/testify/mock"
import sink "github.com/rafalmnich/exporter/sink"

// Exporter is an autogenerated mock type for the Exporter type
type Exporter struct {
	mock.Mock
}

// Export provides a mock function with given fields: ctx, imp
func (_m *Exporter) Export(ctx context.Context, imp []*sink.Reading) error {
	ret := _m.Called(ctx, imp)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*sink.Reading) error); ok {
		r0 = rf(ctx, imp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
