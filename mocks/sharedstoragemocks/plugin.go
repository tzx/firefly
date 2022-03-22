// Code generated by mockery v1.0.0. DO NOT EDIT.

package sharedstoragemocks

import (
	context "context"

	config "github.com/hyperledger/firefly/internal/config"

	io "io"

	mock "github.com/stretchr/testify/mock"

	sharedstorage "github.com/hyperledger/firefly/pkg/sharedstorage"
)

// Plugin is an autogenerated mock type for the Plugin type
type Plugin struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *Plugin) Capabilities() *sharedstorage.Capabilities {
	ret := _m.Called()

	var r0 *sharedstorage.Capabilities
	if rf, ok := ret.Get(0).(func() *sharedstorage.Capabilities); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sharedstorage.Capabilities)
		}
	}

	return r0
}

// DownloadData provides a mock function with given fields: ctx, payloadRef
func (_m *Plugin) DownloadData(ctx context.Context, payloadRef string) (io.ReadCloser, error) {
	ret := _m.Called(ctx, payloadRef)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(context.Context, string) io.ReadCloser); ok {
		r0 = rf(ctx, payloadRef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, payloadRef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields: ctx, prefix, callbacks
func (_m *Plugin) Init(ctx context.Context, prefix config.Prefix, callbacks sharedstorage.Callbacks) error {
	ret := _m.Called(ctx, prefix, callbacks)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, config.Prefix, sharedstorage.Callbacks) error); ok {
		r0 = rf(ctx, prefix, callbacks)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitPrefix provides a mock function with given fields: prefix
func (_m *Plugin) InitPrefix(prefix config.Prefix) {
	_m.Called(prefix)
}

// Name provides a mock function with given fields:
func (_m *Plugin) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// UploadData provides a mock function with given fields: ctx, data
func (_m *Plugin) UploadData(ctx context.Context, data io.Reader) (string, error) {
	ret := _m.Called(ctx, data)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader) string); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, io.Reader) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}