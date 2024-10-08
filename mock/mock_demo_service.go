// Code generated by MockGen. DO NOT EDIT.
// Source: go-app/service (interfaces: DemoService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	schema "go-app/schema"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockDemoService is a mock of DemoService interface.
type MockDemoService struct {
	ctrl     *gomock.Controller
	recorder *MockDemoServiceMockRecorder
}

// MockDemoServiceMockRecorder is the mock recorder for MockDemoService.
type MockDemoServiceMockRecorder struct {
	mock *MockDemoService
}

// NewMockDemoService creates a new mock instance.
func NewMockDemoService(ctrl *gomock.Controller) *MockDemoService {
	mock := &MockDemoService{ctrl: ctrl}
	mock.recorder = &MockDemoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDemoService) EXPECT() *MockDemoServiceMockRecorder {
	return m.recorder
}

// Account_Create mocks base method.
func (m *MockDemoService) Account_Create(arg0 context.Context, arg1 *schema.Account_CreateOpts) (*schema.Account_CreateResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Account_Create", arg0, arg1)
	ret0, _ := ret[0].(*schema.Account_CreateResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Account_Create indicates an expected call of Account_Create.
func (mr *MockDemoServiceMockRecorder) Account_Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Account_Create", reflect.TypeOf((*MockDemoService)(nil).Account_Create), arg0, arg1)
}

// CallAPIForMock mocks base method.
func (m *MockDemoService) CallAPIForMock(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CallAPIForMock", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CallAPIForMock indicates an expected call of CallAPIForMock.
func (mr *MockDemoServiceMockRecorder) CallAPIForMock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallAPIForMock", reflect.TypeOf((*MockDemoService)(nil).CallAPIForMock), arg0, arg1)
}

// DemoFunc mocks base method.
func (m *MockDemoService) DemoFunc(arg0 context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DemoFunc", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// DemoFunc indicates an expected call of DemoFunc.
func (mr *MockDemoServiceMockRecorder) DemoFunc(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DemoFunc", reflect.TypeOf((*MockDemoService)(nil).DemoFunc), arg0)
}

// GetAccountDetailWithTransactions mocks base method.
func (m *MockDemoService) GetAccountDetailWithTransactions(arg0 context.Context, arg1 *schema.AccountTransaction_GetOpts) (*schema.Account_Get, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountDetailWithTransactions", arg0, arg1)
	ret0, _ := ret[0].(*schema.Account_Get)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountDetailWithTransactions indicates an expected call of GetAccountDetailWithTransactions.
func (mr *MockDemoServiceMockRecorder) GetAccountDetailWithTransactions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountDetailWithTransactions", reflect.TypeOf((*MockDemoService)(nil).GetAccountDetailWithTransactions), arg0, arg1)
}

// InsertOne mocks base method.
func (m *MockDemoService) InsertOne(arg0 context.Context, arg1 *schema.InsertOneOpts) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOne", arg0, arg1)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertOne indicates an expected call of InsertOne.
func (mr *MockDemoServiceMockRecorder) InsertOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOne", reflect.TypeOf((*MockDemoService)(nil).InsertOne), arg0, arg1)
}

// SentryDemoFunc mocks base method.
func (m *MockDemoService) SentryDemoFunc(arg0 context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SentryDemoFunc", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// SentryDemoFunc indicates an expected call of SentryDemoFunc.
func (mr *MockDemoServiceMockRecorder) SentryDemoFunc(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SentryDemoFunc", reflect.TypeOf((*MockDemoService)(nil).SentryDemoFunc), arg0)
}

// Transaction_Create mocks base method.
func (m *MockDemoService) Transaction_Create(arg0 context.Context, arg1 *schema.Transaction_CreateOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction_Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction_Create indicates an expected call of Transaction_Create.
func (mr *MockDemoServiceMockRecorder) Transaction_Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction_Create", reflect.TypeOf((*MockDemoService)(nil).Transaction_Create), arg0, arg1)
}
