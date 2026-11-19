package Mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	Domain "task_manager/Domain"
)

// NotifierMock is a testify mock of Usecases.Notifier.
type NotifierMock struct {
	mock.Mock
}

func (_m *NotifierMock) Notify(ctx context.Context, task Domain.Task) error {
	ret := _m.Called(ctx, task)
	return ret.Error(0)
}
