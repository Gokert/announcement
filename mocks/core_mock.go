// Code generated by MockGen. DO NOT EDIT.
// Source: core.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "filmoteka/pkg/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockICore is a mock of ICore interface.
type MockICore struct {
	ctrl     *gomock.Controller
	recorder *MockICoreMockRecorder
}

// MockICoreMockRecorder is the mock recorder for MockICore.
type MockICoreMockRecorder struct {
	mock *MockICore
}

// NewMockICore creates a new mock instance.
func NewMockICore(ctrl *gomock.Controller) *MockICore {
	mock := &MockICore{ctrl: ctrl}
	mock.recorder = &MockICoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICore) EXPECT() *MockICoreMockRecorder {
	return m.recorder
}

// CreateAnnouncement mocks base method.
func (m *MockICore) CreateAnnouncement(announcement *models.Announcement, userId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAnnouncement", announcement, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAnnouncement indicates an expected call of CreateAnnouncement.
func (mr *MockICoreMockRecorder) CreateAnnouncement(announcement, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAnnouncement", reflect.TypeOf((*MockICore)(nil).CreateAnnouncement), announcement, userId)
}

// CreateSession mocks base method.
func (m *MockICore) CreateSession(ctx context.Context, login string) (models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, login)
	ret0, _ := ret[0].(models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockICoreMockRecorder) CreateSession(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockICore)(nil).CreateSession), ctx, login)
}

// CreateUserAccount mocks base method.
func (m *MockICore) CreateUserAccount(login, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserAccount", login, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserAccount indicates an expected call of CreateUserAccount.
func (mr *MockICoreMockRecorder) CreateUserAccount(login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserAccount", reflect.TypeOf((*MockICore)(nil).CreateUserAccount), login, password)
}

// FindActiveSession mocks base method.
func (m *MockICore) FindActiveSession(ctx context.Context, sid string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindActiveSession", ctx, sid)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindActiveSession indicates an expected call of FindActiveSession.
func (mr *MockICoreMockRecorder) FindActiveSession(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindActiveSession", reflect.TypeOf((*MockICore)(nil).FindActiveSession), ctx, sid)
}

// FindUserAccount mocks base method.
func (m *MockICore) FindUserAccount(login, password string) (*models.UserItem, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserAccount", login, password)
	ret0, _ := ret[0].(*models.UserItem)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindUserAccount indicates an expected call of FindUserAccount.
func (mr *MockICoreMockRecorder) FindUserAccount(login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserAccount", reflect.TypeOf((*MockICore)(nil).FindUserAccount), login, password)
}

// FindUserByLogin mocks base method.
func (m *MockICore) FindUserByLogin(login string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByLogin", login)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByLogin indicates an expected call of FindUserByLogin.
func (mr *MockICoreMockRecorder) FindUserByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByLogin", reflect.TypeOf((*MockICore)(nil).FindUserByLogin), login)
}

// GetAnnouncement mocks base method.
func (m *MockICore) GetAnnouncement(id uint64) (*models.Announcement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncement", id)
	ret0, _ := ret[0].(*models.Announcement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncement indicates an expected call of GetAnnouncement.
func (mr *MockICoreMockRecorder) GetAnnouncement(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncement", reflect.TypeOf((*MockICore)(nil).GetAnnouncement), id)
}

// GetAnnouncements mocks base method.
func (m *MockICore) GetAnnouncements(page, pageSize uint64) ([]models.Announcement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncements", page, pageSize)
	ret0, _ := ret[0].([]models.Announcement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncements indicates an expected call of GetAnnouncements.
func (mr *MockICoreMockRecorder) GetAnnouncements(page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncements", reflect.TypeOf((*MockICore)(nil).GetAnnouncements), page, pageSize)
}

// GetUserId mocks base method.
func (m *MockICore) GetUserId(ctx context.Context, sid string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserId", ctx, sid)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserId indicates an expected call of GetUserId.
func (mr *MockICoreMockRecorder) GetUserId(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserId", reflect.TypeOf((*MockICore)(nil).GetUserId), ctx, sid)
}

// GetUserName mocks base method.
func (m *MockICore) GetUserName(ctx context.Context, sid string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserName", ctx, sid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserName indicates an expected call of GetUserName.
func (mr *MockICoreMockRecorder) GetUserName(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserName", reflect.TypeOf((*MockICore)(nil).GetUserName), ctx, sid)
}

// KillSession mocks base method.
func (m *MockICore) KillSession(ctx context.Context, sid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KillSession", ctx, sid)
	ret0, _ := ret[0].(error)
	return ret0
}

// KillSession indicates an expected call of KillSession.
func (mr *MockICoreMockRecorder) KillSession(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KillSession", reflect.TypeOf((*MockICore)(nil).KillSession), ctx, sid)
}

// SearchAnnouncements mocks base method.
func (m *MockICore) SearchAnnouncements(page, pageSize, minCost, maxCost uint64, order string) ([]models.Announcement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAnnouncements", page, pageSize, minCost, maxCost, order)
	ret0, _ := ret[0].([]models.Announcement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchAnnouncements indicates an expected call of SearchAnnouncements.
func (mr *MockICoreMockRecorder) SearchAnnouncements(page, pageSize, minCost, maxCost, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAnnouncements", reflect.TypeOf((*MockICore)(nil).SearchAnnouncements), page, pageSize, minCost, maxCost, order)
}
