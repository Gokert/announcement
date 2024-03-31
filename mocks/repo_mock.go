// Code generated by MockGen. DO NOT EDIT.
// Source: repo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "filmoteka/pkg/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIRepository is a mock of IRepository interface.
type MockIRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRepositoryMockRecorder
}

// MockIRepositoryMockRecorder is the mock recorder for MockIRepository.
type MockIRepositoryMockRecorder struct {
	mock *MockIRepository
}

// NewMockIRepository creates a new mock instance.
func NewMockIRepository(ctrl *gomock.Controller) *MockIRepository {
	mock := &MockIRepository{ctrl: ctrl}
	mock.recorder = &MockIRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRepository) EXPECT() *MockIRepositoryMockRecorder {
	return m.recorder
}

// CreateAnnouncement mocks base method.
func (m *MockIRepository) CreateAnnouncement(announcement *models.Announcement, userId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAnnouncement", announcement, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAnnouncement indicates an expected call of CreateAnnouncement.
func (mr *MockIRepositoryMockRecorder) CreateAnnouncement(announcement, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAnnouncement", reflect.TypeOf((*MockIRepository)(nil).CreateAnnouncement), announcement, userId)
}

// CreateUser mocks base method.
func (m *MockIRepository) CreateUser(login string, password []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", login, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockIRepositoryMockRecorder) CreateUser(login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockIRepository)(nil).CreateUser), login, password)
}

// FindUser mocks base method.
func (m *MockIRepository) FindUser(login string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", login)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockIRepositoryMockRecorder) FindUser(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockIRepository)(nil).FindUser), login)
}

// GetAnnouncement mocks base method.
func (m *MockIRepository) GetAnnouncement(id uint64) (*models.Announcement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncement", id)
	ret0, _ := ret[0].(*models.Announcement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncement indicates an expected call of GetAnnouncement.
func (mr *MockIRepositoryMockRecorder) GetAnnouncement(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncement", reflect.TypeOf((*MockIRepository)(nil).GetAnnouncement), id)
}

// GetAnnouncements mocks base method.
func (m *MockIRepository) GetAnnouncements(page, pageSize uint64) ([]models.Announcement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnouncements", page, pageSize)
	ret0, _ := ret[0].([]models.Announcement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnouncements indicates an expected call of GetAnnouncements.
func (mr *MockIRepositoryMockRecorder) GetAnnouncements(page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnouncements", reflect.TypeOf((*MockIRepository)(nil).GetAnnouncements), page, pageSize)
}

// GetUser mocks base method.
func (m *MockIRepository) GetUser(login string, password []byte) (*models.UserItem, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", login, password)
	ret0, _ := ret[0].(*models.UserItem)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUser indicates an expected call of GetUser.
func (mr *MockIRepositoryMockRecorder) GetUser(login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockIRepository)(nil).GetUser), login, password)
}

// GetUserId mocks base method.
func (m *MockIRepository) GetUserId(login string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserId", login)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserId indicates an expected call of GetUserId.
func (mr *MockIRepositoryMockRecorder) GetUserId(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserId", reflect.TypeOf((*MockIRepository)(nil).GetUserId), login)
}

// SearchAnnouncements mocks base method.
func (m *MockIRepository) SearchAnnouncements(page, pageSize, minCost, maxCost uint64, order string) ([]models.Announcement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAnnouncements", page, pageSize, minCost, maxCost, order)
	ret0, _ := ret[0].([]models.Announcement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchAnnouncements indicates an expected call of SearchAnnouncements.
func (mr *MockIRepositoryMockRecorder) SearchAnnouncements(page, pageSize, minCost, maxCost, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAnnouncements", reflect.TypeOf((*MockIRepository)(nil).SearchAnnouncements), page, pageSize, minCost, maxCost, order)
}
