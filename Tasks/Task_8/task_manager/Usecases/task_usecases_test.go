package Usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"task_manager/Domain"
	"task_manager/Mocks"
	"task_manager/Usecases"
)

type TaskUsecaseSuite struct {
	suite.Suite
	repo *Mocks.TaskRepositoryMock
	uc   Usecases.TaskUsecase
	ctx  context.Context
}

// SetupTest runs before every test method and gives us a fresh mock.
// The table-driven tests below also call it at the start of each case.
func (s *TaskUsecaseSuite) SetupTest() {
	s.repo = new(Mocks.TaskRepositoryMock)
	s.uc = Usecases.NewTaskUsecase(s.repo)
	s.ctx = context.Background()
}

func validTask() Domain.Task {
	return Domain.Task{
		Title:          "Write tests",
		DueDate:        time.Date(2026, 11, 20, 9, 0, 0, 0, time.UTC),
		Status:         "pending",
		ReminderStatus: Domain.ReminderPending,
	}
}

func (s *TaskUsecaseSuite) TestGetAllTasks() {
	repoErr := errors.New("db down")
	tasks := []Domain.Task{validTask()}

	tests := []struct {
		name    string
		setup   func()
		want    []Domain.Task
		wantErr error
	}{
		{"success", func() { s.repo.On("GetAll", s.ctx).Return(tasks, nil).Once() }, tasks, nil},
		{"repository error is propagated", func() { s.repo.On("GetAll", s.ctx).Return(nil, repoErr).Once() }, nil, repoErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, err := s.uc.GetAllTasks(s.ctx)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
				s.Equal(tc.want, got)
			}
			s.repo.AssertExpectations(s.T())
		})
	}
}

func (s *TaskUsecaseSuite) TestGetTaskByID() {
	task := validTask()

	tests := []struct {
		name    string
		setup   func()
		want    Domain.Task
		wantErr error
	}{
		{"success", func() { s.repo.On("GetByID", s.ctx, "t1").Return(task, nil).Once() }, task, nil},
		{"task not found", func() { s.repo.On("GetByID", s.ctx, "t1").Return(Domain.Task{}, Domain.ErrNotFound).Once() }, Domain.Task{}, Domain.ErrNotFound},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, err := s.uc.GetTaskByID(s.ctx, "t1")

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
				s.Equal(tc.want, got)
			}
			s.repo.AssertExpectations(s.T())
		})
	}
}

func (s *TaskUsecaseSuite) TestCreateTask() {
	valid := validTask()
	repoErr := errors.New("insert failed")

	blankTitle := validTask()
	blankTitle.Title = "   "
	noDate := validTask()
	noDate.DueDate = time.Time{}
	badStatus := validTask()
	badStatus.Status = "finished"

	tests := []struct {
		name    string
		input   Domain.Task
		setup   func()
		wantErr error
	}{
		{"success", valid, func() { s.repo.On("Create", s.ctx, valid).Return(valid, nil).Once() }, nil},
		{"blank title is rejected, repo never called", blankTitle, func() {}, Domain.ErrInvalidInput},
		{"missing due date is rejected", noDate, func() {}, Domain.ErrInvalidInput},
		{"invalid status is rejected", badStatus, func() {}, Domain.ErrInvalidInput},
		{"repository error is propagated", valid, func() { s.repo.On("Create", s.ctx, valid).Return(Domain.Task{}, repoErr).Once() }, repoErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, err := s.uc.CreateTask(s.ctx, tc.input)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
				s.Equal(tc.input, got)
			}
			// An unexpected repo call would panic, so a passing "rejected" case
			// also proves the repository was never touched.
			s.repo.AssertExpectations(s.T())
		})
	}
}

func (s *TaskUsecaseSuite) TestUpdateTask() {
	valid := validTask()
	repoErr := errors.New("update failed")
	blankTitle := validTask()
	blankTitle.Title = ""

	tests := []struct {
		name    string
		input   Domain.Task
		setup   func()
		wantErr error
	}{
		{"success", valid, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(valid, nil).Once()
			s.repo.On("Update", s.ctx, "t1", valid).Return(valid, nil).Once()
		}, nil},
		{"validation fails before any repo call", blankTitle, func() {}, Domain.ErrInvalidInput},
		{"task not found, Update never called", valid, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(Domain.Task{}, Domain.ErrNotFound).Once()
		}, Domain.ErrNotFound},
		{"repository update error is propagated", valid, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(valid, nil).Once()
			s.repo.On("Update", s.ctx, "t1", valid).Return(Domain.Task{}, repoErr).Once()
		}, repoErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, err := s.uc.UpdateTask(s.ctx, "t1", tc.input)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
				s.Equal(tc.input, got)
			}
			s.repo.AssertExpectations(s.T())
		})
	}
}

func (s *TaskUsecaseSuite) TestDeleteTask() {
	repoErr := errors.New("delete failed")

	tests := []struct {
		name    string
		setup   func()
		wantErr error
	}{
		{"success", func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(validTask(), nil).Once()
			s.repo.On("Delete", s.ctx, "t1").Return(nil).Once()
		}, nil},
		{"task not found, Delete never called", func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(Domain.Task{}, Domain.ErrNotFound).Once()
		}, Domain.ErrNotFound},
		{"repository delete error is propagated", func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(validTask(), nil).Once()
			s.repo.On("Delete", s.ctx, "t1").Return(repoErr).Once()
		}, repoErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			err := s.uc.DeleteTask(s.ctx, "t1")

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
			}
			s.repo.AssertExpectations(s.T())
		})
	}
}

func (s *TaskUsecaseSuite) TestCreateTaskAlwaysStartsPending() {
	in := validTask()
	in.ReminderStatus = Domain.ReminderSent // a client trying to set it
	want := in
	want.ReminderStatus = Domain.ReminderPending
	s.repo.On("Create", s.ctx, want).Return(want, nil).Once()

	got, err := s.uc.CreateTask(s.ctx, in)

	s.NoError(err)
	s.Equal(Domain.ReminderPending, got.ReminderStatus)
	s.repo.AssertExpectations(s.T())
}

func (s *TaskUsecaseSuite) TestUpdateTaskReminderStatus() {
	existing := validTask()
	existing.ReminderStatus = Domain.ReminderSent

	tests := []struct {
		name       string
		newDue     time.Time
		wantStatus string
	}{
		{"due date moved to the future resets to pending", time.Now().Add(48 * time.Hour), Domain.ReminderPending},
		{"due date unchanged keeps the stored status", existing.DueDate, Domain.ReminderSent},
		{"due date moved to the past keeps the stored status", time.Now().Add(-time.Hour), Domain.ReminderSent},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			in := existing
			in.DueDate = tc.newDue
			in.ReminderStatus = "" // client value is ignored
			want := in
			want.ReminderStatus = tc.wantStatus
			s.repo.On("GetByID", s.ctx, "t1").Return(existing, nil).Once()
			s.repo.On("Update", s.ctx, "t1", want).Return(want, nil).Once()

			got, err := s.uc.UpdateTask(s.ctx, "t1", in)

			s.NoError(err)
			s.Equal(tc.wantStatus, got.ReminderStatus)
			s.repo.AssertExpectations(s.T())
		})
	}
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseSuite))
}
