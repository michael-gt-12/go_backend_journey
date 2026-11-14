package Usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"task_manager/Domain"
	"task_manager/Mocks"
	"task_manager/Usecases"
)

type ReminderUsecaseSuite struct {
	suite.Suite
	repo     *Mocks.TaskRepositoryMock
	notifier *Mocks.NotifierMock
}

func (s *ReminderUsecaseSuite) SetupTest() {
	s.repo = new(Mocks.TaskRepositoryMock)
	s.notifier = new(Mocks.NotifierMock)
}

func reminderTask() Domain.Task {
	return Domain.Task{
		ID:             "t1",
		Title:          "Submit report",
		DueDate:        time.Date(2026, 11, 20, 9, 0, 0, 0, time.UTC),
		Status:         "pending",
		ReminderStatus: Domain.ReminderQueued,
	}
}

func (s *ReminderUsecaseSuite) TestDeliver() {
	task := reminderTask()
	errBoom := errors.New("notifier down")
	errDB := errors.New("db down")

	tests := []struct {
		name       string
		attempts   int
		retryDelay time.Duration
		setup      func(cancel context.CancelFunc)
		wantErr    error
	}{
		{"success on the first try", 3, time.Millisecond, func(context.CancelFunc) {
			s.notifier.On("Notify", mock.Anything, task).Return(nil).Once()
			s.repo.On("SetReminderStatus", mock.Anything, "t1", Domain.ReminderSent).Return(nil).Once()
		}, nil},
		{"success on the second try", 3, time.Millisecond, func(context.CancelFunc) {
			s.notifier.On("Notify", mock.Anything, task).Return(errBoom).Once()
			s.notifier.On("Notify", mock.Anything, task).Return(nil).Once()
			s.repo.On("SetReminderStatus", mock.Anything, "t1", Domain.ReminderSent).Return(nil).Once()
		}, nil},
		{"failure after all attempts marks failed", 3, time.Millisecond, func(context.CancelFunc) {
			s.notifier.On("Notify", mock.Anything, task).Return(errBoom).Times(3)
			s.repo.On("SetReminderStatus", mock.Anything, "t1", Domain.ReminderFailed).Return(nil).Once()
		}, errBoom},
		{"ctx cancelled during the wait stops immediately", 3, time.Hour, func(cancel context.CancelFunc) {
			// Cancel from inside the first Notify; the one-hour wait must be skipped.
			s.notifier.On("Notify", mock.Anything, task).
				Run(func(mock.Arguments) { cancel() }).Return(errBoom).Once()
			// No SetReminderStatus expectation: calling it would panic the mock.
		}, context.Canceled},
		{"saving the sent status fails", 3, time.Millisecond, func(context.CancelFunc) {
			s.notifier.On("Notify", mock.Anything, task).Return(nil).Once()
			s.repo.On("SetReminderStatus", mock.Anything, "t1", Domain.ReminderSent).Return(errDB).Once()
		}, errDB},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			tc.setup(cancel)
			uc := Usecases.NewReminderUsecase(s.repo, s.notifier, tc.attempts, tc.retryDelay)

			err := uc.Deliver(ctx, task)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
			}
			s.notifier.AssertExpectations(s.T())
			s.repo.AssertExpectations(s.T())
		})
	}
}

func (s *ReminderUsecaseSuite) TestClaimDue() {
	window := 10 * time.Minute
	claimed := []Domain.Task{reminderTask()}
	repoErr := errors.New("claim failed")

	tests := []struct {
		name    string
		ret     []Domain.Task
		retErr  error
		wantErr error
	}{
		{"returns the claimed tasks", claimed, nil, nil},
		{"repository error is propagated", nil, repoErr, repoErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			before := time.Now()
			// dueBefore must be "now + window": check it falls between two clock readings.
			dueOK := mock.MatchedBy(func(t time.Time) bool {
				return !t.Before(before.Add(window)) && !t.After(time.Now().Add(window))
			})
			s.repo.On("ClaimDueTasks", mock.Anything, dueOK, mock.Anything).Return(tc.ret, tc.retErr).Once()
			uc := Usecases.NewReminderUsecase(s.repo, s.notifier, 3, time.Millisecond)

			got, err := uc.ClaimDue(context.Background(), window)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
				s.Equal(tc.ret, got)
			}
			s.repo.AssertExpectations(s.T())
		})
	}
}

func TestReminderUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ReminderUsecaseSuite))
}
