// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/user/dip_data.go

// Package user is a generated GoMock package.
package user

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	sessions "yumi/pkg/sessions"
	entity "yumi/usecase/user/entity"
)

// MockData is a mock of Data interface
type MockData struct {
	ctrl     *gomock.Controller
	recorder *MockDataMockRecorder
}

// MockDataMockRecorder is the mock recorder for MockData
type MockDataMockRecorder struct {
	mock *MockData
}

// NewMockData creates a new mock instance
func NewMockData(ctrl *gomock.Controller) *MockData {
	mock := &MockData{ctrl: ctrl}
	mock.recorder = &MockDataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockData) EXPECT() *MockDataMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockData) Create(arg0 entity.UserAttribute) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockDataMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockData)(nil).Create), arg0)
}

// GetUser mocks base method
func (m *MockData) GetUser(userID string) (DataUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", userID)
	ret0, _ := ret[0].(DataUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockDataMockRecorder) GetUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockData)(nil).GetUser), userID)
}

// GetSessionsStore mocks base method
func (m *MockData) GetSessionsStore() sessions.Store {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionsStore")
	ret0, _ := ret[0].(sessions.Store)
	return ret0
}

// GetSessionsStore indicates an expected call of GetSessionsStore
func (mr *MockDataMockRecorder) GetSessionsStore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionsStore", reflect.TypeOf((*MockData)(nil).GetSessionsStore))
}

// MockDataUser is a mock of DataUser interface
type MockDataUser struct {
	ctrl     *gomock.Controller
	recorder *MockDataUserMockRecorder
}

// MockDataUserMockRecorder is the mock recorder for MockDataUser
type MockDataUserMockRecorder struct {
	mock *MockDataUser
}

// NewMockDataUser creates a new mock instance
func NewMockDataUser(ctrl *gomock.Controller) *MockDataUser {
	mock := &MockDataUser{ctrl: ctrl}
	mock.recorder = &MockDataUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataUser) EXPECT() *MockDataUserMockRecorder {
	return m.recorder
}

// Attribute mocks base method
func (m *MockDataUser) Attribute() *entity.UserAttribute {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Attribute")
	ret0, _ := ret[0].(*entity.UserAttribute)
	return ret0
}

// Attribute indicates an expected call of Attribute
func (mr *MockDataUserMockRecorder) Attribute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attribute", reflect.TypeOf((*MockDataUser)(nil).Attribute))
}

// Update mocks base method
func (m *MockDataUser) Update() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update")
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockDataUserMockRecorder) Update() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDataUser)(nil).Update))
}