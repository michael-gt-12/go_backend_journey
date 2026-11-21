package Usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"task_manager/Domain"
	"task_manager/Mocks"
	"task_manager/Repositories"
	"task_manager/Usecases"
)

type TaskUsecaseSuite struct {
	suite.Suite
	repo     *Mocks.TaskRepositoryMock
	uc       Usecases.TaskUsecase
	ctx      context.Context // normal user "u1"
	adminCtx context.Context // admin "a1"
}

// callerCtx builds a context carrying an authenticated caller.
func callerCtx(id, role string) context.Context {
	return Usecases.WithCaller(context.Background(), Usecases.Caller{ID: id, Role: role})
}

// SetupTest runs before every test method and gives us a fresh mock.
// The table-driven tests below also call it at the start of each case.
func (s *TaskUsecaseSuite) SetupTest() {
	s.repo = new(Mocks.TaskRepositoryMock)
	s.uc = Usecases.NewTaskUsecase(s.repo)
	s.ctx = callerCtx("u1", "user")
	s.adminCtx = callerCtx("a1", "admin")
}

// validTask is a task owned by the normal user "u1".
func validTask() Domain.Task {
	return Domain.Task{
		OwnerID:        "u1",
		Title:          "Write tests",
		DueDate:        time.Date(2026, 11, 20, 9, 0, 0, 0, time.UTC),
		Status:         "pending",
		ReminderStatus: Domain.ReminderPending,
	}
}

// otherUsersTask is a task owned by someone else ("u2").
func otherUsersTask() Domain.Task {
	t := validTask()
	t.OwnerID = "u2"
	return t
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

func (s *TaskUsecaseSuite) TestListTasks() {
	repoErr := errors.New("db down")
	tasks := []Domain.Task{validTask()}

	tests := []struct {
		name      string
		ctx       func() context.Context
		filter    Repositories.TaskFilter
		setup     func()
		wantTotal int
		wantErr   error
	}{
		{
			name:   "normal user is forced to their own tasks and defaults are applied",
			ctx:    func() context.Context { return s.ctx },
			filter: Repositories.TaskFilter{},
			setup: func() {
				s.repo.On("List", s.ctx, Repositories.TaskFilter{OwnerID: "u1", Page: 1, Limit: 10}).
					Return(tasks, 1, nil).Once()
			},
			wantTotal: 1,
		},
		{
			name:   "normal user cannot ask for another owner's tasks",
			ctx:    func() context.Context { return s.ctx },
			filter: Repositories.TaskFilter{OwnerID: "u2", Page: 2, Limit: 5},
			setup: func() {
				s.repo.On("List", s.ctx, Repositories.TaskFilter{OwnerID: "u1", Page: 2, Limit: 5}).
					Return(tasks, 7, nil).Once()
			},
			wantTotal: 7,
		},
		{
			name:   "admin sees all owners (empty OwnerID)",
			ctx:    func() context.Context { return s.adminCtx },
			filter: Repositories.TaskFilter{Status: "done"},
			setup: func() {
				s.repo.On("List", s.adminCtx, Repositories.TaskFilter{Status: "done", Page: 1, Limit: 10}).
					Return(tasks, 1, nil).Once()
			},
			wantTotal: 1,
		},
		{
			name:    "invalid status is rejected, repo never called",
			ctx:     func() context.Context { return s.ctx },
			filter:  Repositories.TaskFilter{Status: "finished"},
			setup:   func() {},
			wantErr: Domain.ErrInvalidInput,
		},
		{
			name:    "missing caller is rejected, repo never called",
			ctx:     func() context.Context { return context.Background() },
			filter:  Repositories.TaskFilter{},
			setup:   func() {},
			wantErr: Domain.ErrInvalidInput,
		},
		{
			name:   "repository error is propagated",
			ctx:    func() context.Context { return s.ctx },
			filter: Repositories.TaskFilter{},
			setup: func() {
				s.repo.On("List", s.ctx, Repositories.TaskFilter{OwnerID: "u1", Page: 1, Limit: 10}).
					Return(nil, 0, repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, total, err := s.uc.ListTasks(tc.ctx(), tc.filter)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
			} else {
				s.NoError(err)
				s.Equal(tasks, got)
				s.Equal(tc.wantTotal, total)
			}
			s.repo.AssertExpectations(s.T())
		})
	}
}

func (s *TaskUsecaseSuite) TestGetTaskByID() {
	tests := []struct {
		name    string
		ctx     func() context.Context
		setup   func()
		want    Domain.Task
		wantErr error
	}{
		{"owner can read their task", func() context.Context { return s.ctx },
			func() { s.repo.On("GetByID", s.ctx, "t1").Return(validTask(), nil).Once() },
			validTask(), nil},
		{"another user's task looks like 404", func() context.Context { return s.ctx },
			func() { s.repo.On("GetByID", s.ctx, "t1").Return(otherUsersTask(), nil).Once() },
			Domain.Task{}, Domain.ErrNotFound},
		{"admin can read anyone's task", func() context.Context { return s.adminCtx },
			func() { s.repo.On("GetByID", s.adminCtx, "t1").Return(otherUsersTask(), nil).Once() },
			otherUsersTask(), nil},
		{"task not found", func() context.Context { return s.ctx },
			func() { s.repo.On("GetByID", s.ctx, "t1").Return(Domain.Task{}, Domain.ErrNotFound).Once() },
			Domain.Task{}, Domain.ErrNotFound},
		{"missing caller is rejected, repo never called", func() context.Context { return context.Background() },
			func() {}, Domain.Task{}, Domain.ErrInvalidInput},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, err := s.uc.GetTaskByID(tc.ctx(), "t1")

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
	valid := validTask() // OwnerID u1, the caller
	repoErr := errors.New("insert failed")

	blankTitle := validTask()
	blankTitle.Title = "   "
	noDate := validTask()
	noDate.DueDate = time.Time{}
	badStatus := validTask()
	badStatus.Status = "finished"

	tests := []struct {
		name    string
		ctx     func() context.Context
		input   Domain.Task
		setup   func()
		wantErr error
	}{
		{"success", func() context.Context { return s.ctx }, valid,
			func() { s.repo.On("Create", s.ctx, valid).Return(valid, nil).Once() }, nil},
		{"blank title is rejected, repo never called", func() context.Context { return s.ctx }, blankTitle, func() {}, Domain.ErrInvalidInput},
		{"missing due date is rejected", func() context.Context { return s.ctx }, noDate, func() {}, Domain.ErrInvalidInput},
		{"invalid status is rejected", func() context.Context { return s.ctx }, badStatus, func() {}, Domain.ErrInvalidInput},
		{"missing caller is rejected", func() context.Context { return context.Background() }, valid, func() {}, Domain.ErrInvalidInput},
		{"repository error is propagated", func() context.Context { return s.ctx }, valid,
			func() { s.repo.On("Create", s.ctx, valid).Return(Domain.Task{}, repoErr).Once() }, repoErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, err := s.uc.CreateTask(tc.ctx(), tc.input)

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

// The owner always comes from the authenticated caller, never from the request body.
func (s *TaskUsecaseSuite) TestCreateTaskOwnerComesFromCaller() {
	in := validTask()
	in.OwnerID = "evil-user" // a client trying to create a task for someone else
	want := validTask()      // OwnerID u1
	s.repo.On("Create", s.ctx, want).Return(want, nil).Once()

	got, err := s.uc.CreateTask(s.ctx, in)

	s.NoError(err)
	s.Equal("u1", got.OwnerID)
	s.repo.AssertExpectations(s.T())
}

func (s *TaskUsecaseSuite) TestUpdateTask() {
	valid := validTask()
	repoErr := errors.New("update failed")
	blankTitle := validTask()
	blankTitle.Title = ""

	// An admin updating u2's task must keep u2 as the owner.
	adminInput := validTask()
	adminInput.OwnerID = ""
	adminWant := otherUsersTask()

	tests := []struct {
		name    string
		ctx     func() context.Context
		input   Domain.Task
		setup   func()
		want    Domain.Task
		wantErr error
	}{
		{"owner updates their task", func() context.Context { return s.ctx }, valid, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(valid, nil).Once()
			s.repo.On("Update", s.ctx, "t1", valid).Return(valid, nil).Once()
		}, valid, nil},
		{"admin updates another user's task and the owner is preserved", func() context.Context { return s.adminCtx }, adminInput, func() {
			s.repo.On("GetByID", s.adminCtx, "t1").Return(otherUsersTask(), nil).Once()
			s.repo.On("Update", s.adminCtx, "t1", adminWant).Return(adminWant, nil).Once()
		}, adminWant, nil},
		{"validation fails before any repo call", func() context.Context { return s.ctx }, blankTitle, func() {}, Domain.Task{}, Domain.ErrInvalidInput},
		{"another user's task looks like 404, Update never called", func() context.Context { return s.ctx }, valid, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(otherUsersTask(), nil).Once()
		}, Domain.Task{}, Domain.ErrNotFound},
		{"task not found, Update never called", func() context.Context { return s.ctx }, valid, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(Domain.Task{}, Domain.ErrNotFound).Once()
		}, Domain.Task{}, Domain.ErrNotFound},
		{"missing caller is rejected", func() context.Context { return context.Background() }, valid, func() {}, Domain.Task{}, Domain.ErrInvalidInput},
		{"repository update error is propagated", func() context.Context { return s.ctx }, valid, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(valid, nil).Once()
			s.repo.On("Update", s.ctx, "t1", valid).Return(Domain.Task{}, repoErr).Once()
		}, Domain.Task{}, repoErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			got, err := s.uc.UpdateTask(tc.ctx(), "t1", tc.input)

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

func (s *TaskUsecaseSuite) TestDeleteTask() {
	repoErr := errors.New("delete failed")

	tests := []struct {
		name    string
		ctx     func() context.Context
		setup   func()
		wantErr error
	}{
		{"owner deletes their task", func() context.Context { return s.ctx }, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(validTask(), nil).Once()
			s.repo.On("Delete", s.ctx, "t1").Return(nil).Once()
		}, nil},
		{"admin deletes another user's task", func() context.Context { return s.adminCtx }, func() {
			s.repo.On("GetByID", s.adminCtx, "t1").Return(otherUsersTask(), nil).Once()
			s.repo.On("Delete", s.adminCtx, "t1").Return(nil).Once()
		}, nil},
		{"another user's task looks like 404, Delete never called", func() context.Context { return s.ctx }, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(otherUsersTask(), nil).Once()
		}, Domain.ErrNotFound},
		{"task not found, Delete never called", func() context.Context { return s.ctx }, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(Domain.Task{}, Domain.ErrNotFound).Once()
		}, Domain.ErrNotFound},
		{"missing caller is rejected", func() context.Context { return context.Background() }, func() {}, Domain.ErrInvalidInput},
		{"repository delete error is propagated", func() context.Context { return s.ctx }, func() {
			s.repo.On("GetByID", s.ctx, "t1").Return(validTask(), nil).Once()
			s.repo.On("Delete", s.ctx, "t1").Return(repoErr).Once()
		}, repoErr},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()

			err := s.uc.DeleteTask(tc.ctx(), "t1")

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
	existing := validTask() // owned by u1
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
			in.OwnerID = ""        // client value is ignored
			in.ReminderStatus = "" // client value is ignored
			want := in
			want.OwnerID = "u1" // kept from the stored task
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
