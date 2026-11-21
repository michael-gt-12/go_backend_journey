package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"task_manager/Delivery/controllers"
	"task_manager/Delivery/routers"
	"task_manager/Domain"
	"task_manager/Infrastructure"
	"task_manager/Mocks"
	"task_manager/Repositories"
)

type ControllerSuite struct {
	suite.Suite
	tasks      *Mocks.TaskUsecaseMock
	users      *Mocks.UserUsecaseMock
	router     *gin.Engine
	adminToken string
	userToken  string
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}

// SetupTest runs before every test: fresh mocks and a fresh router.
func (s *ControllerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.tasks = new(Mocks.TaskUsecaseMock)
	s.users = new(Mocks.UserUsecaseMock)

	jwtService := Infrastructure.NewJWTService("test-secret", time.Hour)
	h := controllers.NewController(s.tasks, s.users)
	s.router = routers.SetupRouter(h, jwtService)

	var err error
	s.adminToken, err = jwtService.Generate(Domain.User{ID: "a1", Username: "admin", Role: Domain.RoleAdmin})
	s.Require().NoError(err)
	s.userToken, err = jwtService.Generate(Domain.User{ID: "u1", Username: "bob", Role: Domain.RoleUser})
	s.Require().NoError(err)
}

// do sends a request through the Gin router. body may be a string (raw) or any value (marshalled to JSON).
func (s *ControllerSuite) do(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	switch b := body.(type) {
	case nil:
	case string:
		buf.WriteString(b)
	default:
		s.Require().NoError(json.NewEncoder(&buf).Encode(b))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	return w
}

type tc struct {
	name     string
	method   string
	path     string
	body     interface{}
	token    func() string // evaluated after SetupTest so tokens exist
	setup    func()        // sets mock expectations (nil = usecase must NOT be called)
	status   int
	contains string
	headers  map[string]string // expected response headers
}

func (s *ControllerSuite) run(cases []tc) {
	for _, c := range cases {
		s.Run(c.name, func() {
			s.SetupTest() // fresh mocks for every case
			if c.setup != nil {
				c.setup()
			}
			token := ""
			if c.token != nil {
				token = c.token()
			}
			w := s.do(c.method, c.path, c.body, token)

			s.Equal(c.status, w.Code, w.Body.String())
			if c.contains != "" {
				s.Contains(w.Body.String(), c.contains)
			}
			for k, v := range c.headers {
				s.Equal(v, w.Header().Get(k), "header "+k)
			}
			s.tasks.AssertExpectations(s.T())
			s.users.AssertExpectations(s.T())
		})
	}
}

func (s *ControllerSuite) admin() string { return s.adminToken }
func (s *ControllerSuite) user() string  { return s.userToken }

var due = time.Date(2026, 12, 1, 10, 0, 0, 0, time.UTC)

func sampleTask() Domain.Task {
	return Domain.Task{ID: "t1", Title: "Write tests", Description: "controllers", DueDate: due, Status: "pending"}
}

// ---------- auth endpoints ----------

func (s *ControllerSuite) TestRegister() {
	s.run([]tc{
		{
			name: "success returns 201", method: "POST", path: "/register",
			body: map[string]string{"username": "bob", "password": "secret1"},
			setup: func() {
				s.users.On("Register", mock.Anything, "bob", "secret1").
					Return(Domain.User{ID: "u1", Username: "bob", Role: Domain.RoleUser}, nil).Once()
			},
			status: http.StatusCreated, contains: "bob",
		},
		{
			name: "missing password returns 400", method: "POST", path: "/register",
			body:   map[string]string{"username": "bob"},
			status: http.StatusBadRequest,
		},
		{
			name: "malformed JSON returns 400", method: "POST", path: "/register",
			body: "{not json", status: http.StatusBadRequest,
		},
		{
			name: "duplicate username returns 409", method: "POST", path: "/register",
			body: map[string]string{"username": "bob", "password": "secret1"},
			setup: func() {
				s.users.On("Register", mock.Anything, "bob", "secret1").
					Return(Domain.User{}, Domain.ErrUsernameTaken).Once()
			},
			status: http.StatusConflict,
		},
		{
			name: "invalid input from usecase returns 400", method: "POST", path: "/register",
			body: map[string]string{"username": "bo", "password": "secret1"},
			setup: func() {
				s.users.On("Register", mock.Anything, "bo", "secret1").
					Return(Domain.User{}, Domain.ErrInvalidInput).Once()
			},
			status: http.StatusBadRequest,
		},
	})
}

func (s *ControllerSuite) TestLogin() {
	s.run([]tc{
		{
			name: "success returns 200 and token", method: "POST", path: "/login",
			body: map[string]string{"username": "bob", "password": "secret1"},
			setup: func() {
				s.users.On("Login", mock.Anything, "bob", "secret1").Return("jwt-token", nil).Once()
			},
			status: http.StatusOK, contains: "jwt-token",
		},
		{
			name: "wrong credentials return 401", method: "POST", path: "/login",
			body: map[string]string{"username": "bob", "password": "wrong"},
			setup: func() {
				s.users.On("Login", mock.Anything, "bob", "wrong").
					Return("", Domain.ErrInvalidCredentials).Once()
			},
			status: http.StatusUnauthorized, contains: "invalid username or password",
		},
		{
			name: "missing fields return 400", method: "POST", path: "/login",
			body: map[string]string{"username": "bob"}, status: http.StatusBadRequest,
		},
		{
			name: "unexpected error returns 500 without leaking details", method: "POST", path: "/login",
			body: map[string]string{"username": "bob", "password": "secret1"},
			setup: func() {
				s.users.On("Login", mock.Anything, "bob", "secret1").
					Return("", errors.New("db exploded")).Once()
			},
			status: http.StatusInternalServerError, contains: "internal server error",
		},
	})
}

func (s *ControllerSuite) TestPromote() {
	s.run([]tc{
		{
			name: "admin can promote", method: "POST", path: "/promote/u1", token: s.admin,
			setup:  func() { s.users.On("Promote", mock.Anything, "u1").Return(nil).Once() },
			status: http.StatusOK,
		},
		{
			name: "normal user gets 403", method: "POST", path: "/promote/u1", token: s.user,
			status: http.StatusForbidden,
		},
		{
			name: "no token gets 401", method: "POST", path: "/promote/u1",
			status: http.StatusUnauthorized,
		},
		{
			name: "unknown user returns 404", method: "POST", path: "/promote/nope", token: s.admin,
			setup:  func() { s.users.On("Promote", mock.Anything, "nope").Return(Domain.ErrNotFound).Once() },
			status: http.StatusNotFound,
		},
	})
}

// ---------- task endpoints ----------

func (s *ControllerSuite) TestGetTasks() {
	defaults := Repositories.TaskFilter{Page: 1, Limit: 10}

	s.run([]tc{
		{
			name: "logged-in user gets 200, default filter and pagination headers", method: "GET", path: "/tasks", token: s.user,
			setup: func() {
				s.tasks.On("ListTasks", mock.Anything, defaults).Return([]Domain.Task{sampleTask()}, 42, nil).Once()
			},
			status: http.StatusOK, contains: "Write tests",
			headers: map[string]string{"X-Total-Count": "42", "X-Page": "1", "X-Limit": "10"},
		},
		{
			name: "query parameters are parsed", method: "GET", path: "/tasks?page=2&limit=5&status=done", token: s.user,
			setup: func() {
				s.tasks.On("ListTasks", mock.Anything, Repositories.TaskFilter{Page: 2, Limit: 5, Status: "done"}).
					Return([]Domain.Task{sampleTask()}, 11, nil).Once()
			},
			status:  http.StatusOK,
			headers: map[string]string{"X-Total-Count": "11", "X-Page": "2", "X-Limit": "5"},
		},
		{
			name: "admin can list too", method: "GET", path: "/tasks", token: s.admin,
			setup: func() {
				s.tasks.On("ListTasks", mock.Anything, defaults).Return([]Domain.Task{sampleTask()}, 1, nil).Once()
			},
			status: http.StatusOK,
		},
		{
			name: "empty result is an empty array", method: "GET", path: "/tasks", token: s.user,
			setup: func() {
				s.tasks.On("ListTasks", mock.Anything, defaults).Return([]Domain.Task{}, 0, nil).Once()
			},
			status: http.StatusOK, contains: "[]",
			headers: map[string]string{"X-Total-Count": "0"},
		},
		{name: "page=0 returns 400", method: "GET", path: "/tasks?page=0", token: s.user, status: http.StatusBadRequest},
		{name: "page=abc returns 400", method: "GET", path: "/tasks?page=abc", token: s.user, status: http.StatusBadRequest},
		{name: "limit=0 returns 400", method: "GET", path: "/tasks?limit=0", token: s.user, status: http.StatusBadRequest},
		{name: "limit=101 returns 400", method: "GET", path: "/tasks?limit=101", token: s.user, status: http.StatusBadRequest},
		{name: "unknown status returns 400", method: "GET", path: "/tasks?status=foo", token: s.user, status: http.StatusBadRequest},
		{
			name: "no token returns 401", method: "GET", path: "/tasks", status: http.StatusUnauthorized,
		},
		{
			name: "tampered token returns 401", method: "GET", path: "/tasks",
			token:  func() string { return "not.a.jwt" },
			status: http.StatusUnauthorized,
		},
		{
			name: "usecase error returns 500", method: "GET", path: "/tasks", token: s.user,
			setup: func() {
				s.tasks.On("ListTasks", mock.Anything, defaults).Return(nil, 0, errors.New("db down")).Once()
			},
			status: http.StatusInternalServerError, contains: "internal server error",
		},
	})
}

func (s *ControllerSuite) TestGetTask() {
	s.run([]tc{
		{
			name: "found returns 200", method: "GET", path: "/tasks/t1", token: s.user,
			setup:  func() { s.tasks.On("GetTaskByID", mock.Anything, "t1").Return(sampleTask(), nil).Once() },
			status: http.StatusOK, contains: "Write tests",
		},
		{
			name: "missing or someone else's task returns 404", method: "GET", path: "/tasks/zzz", token: s.user,
			setup: func() {
				s.tasks.On("GetTaskByID", mock.Anything, "zzz").Return(Domain.Task{}, Domain.ErrNotFound).Once()
			},
			status: http.StatusNotFound,
		},
		{
			name: "no token returns 401", method: "GET", path: "/tasks/t1", status: http.StatusUnauthorized,
		},
	})
}

func (s *ControllerSuite) TestCreateTask() {
	// The controller forces ID = "" and OwnerID = "" before calling the usecase;
	// the usecase then sets the owner from the authenticated caller.
	matchNew := mock.MatchedBy(func(t Domain.Task) bool {
		return t.ID == "" && t.OwnerID == "" && t.Title == "Write tests" && t.Status == "pending" && t.DueDate.Equal(due)
	})
	s.run([]tc{
		{
			name: "admin creates task -> 201", method: "POST", path: "/tasks", token: s.admin,
			body: sampleTask(),
			setup: func() {
				s.tasks.On("CreateTask", mock.Anything, matchNew).Return(sampleTask(), nil).Once()
			},
			status: http.StatusCreated, contains: "Write tests",
		},
		{
			name: "normal user creates their own task -> 201", method: "POST", path: "/tasks", token: s.user,
			body: sampleTask(),
			setup: func() {
				s.tasks.On("CreateTask", mock.Anything, matchNew).Return(sampleTask(), nil).Once()
			},
			status: http.StatusCreated,
		},
		{
			name: "owner_id in the body is ignored", method: "POST", path: "/tasks", token: s.user,
			body: map[string]interface{}{
				"title": "Write tests", "due_date": due, "status": "pending", "owner_id": "evil-user",
			},
			setup: func() {
				s.tasks.On("CreateTask", mock.Anything, matchNew).Return(sampleTask(), nil).Once()
			},
			status: http.StatusCreated,
		},
		{
			name: "invalid JSON -> 400", method: "POST", path: "/tasks", token: s.user,
			body: "{bad json", status: http.StatusBadRequest,
		},
		{
			name: "usecase validation error -> 400", method: "POST", path: "/tasks", token: s.user,
			body: Domain.Task{Description: "no title"},
			setup: func() {
				s.tasks.On("CreateTask", mock.Anything, mock.Anything).
					Return(Domain.Task{}, Domain.ErrInvalidInput).Once()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "no token -> 401", method: "POST", path: "/tasks",
			body: sampleTask(), status: http.StatusUnauthorized,
		},
	})
}

func (s *ControllerSuite) TestUpdateTask() {
	s.run([]tc{
		{
			name: "owner updates task -> 200", method: "PUT", path: "/tasks/t1", token: s.user,
			body: sampleTask(),
			setup: func() {
				s.tasks.On("UpdateTask", mock.Anything, "t1", mock.MatchedBy(func(t Domain.Task) bool {
					return t.Title == "Write tests" && t.OwnerID == ""
				})).Return(sampleTask(), nil).Once()
			},
			status: http.StatusOK,
		},
		{
			name: "admin updates task -> 200", method: "PUT", path: "/tasks/t1", token: s.admin,
			body: sampleTask(),
			setup: func() {
				s.tasks.On("UpdateTask", mock.Anything, "t1", mock.Anything).Return(sampleTask(), nil).Once()
			},
			status: http.StatusOK,
		},
		{
			name: "invalid JSON -> 400", method: "PUT", path: "/tasks/t1", token: s.user,
			body: "{bad", status: http.StatusBadRequest,
		},
		{
			name: "missing or someone else's task -> 404", method: "PUT", path: "/tasks/zzz", token: s.user,
			body: sampleTask(),
			setup: func() {
				s.tasks.On("UpdateTask", mock.Anything, "zzz", mock.Anything).
					Return(Domain.Task{}, Domain.ErrNotFound).Once()
			},
			status: http.StatusNotFound,
		},
		{
			name: "no token -> 401", method: "PUT", path: "/tasks/t1",
			body: sampleTask(), status: http.StatusUnauthorized,
		},
	})
}

func (s *ControllerSuite) TestDeleteTask() {
	s.run([]tc{
		{
			name: "owner deletes task -> 200", method: "DELETE", path: "/tasks/t1", token: s.user,
			setup:  func() { s.tasks.On("DeleteTask", mock.Anything, "t1").Return(nil).Once() },
			status: http.StatusOK, contains: "task deleted",
		},
		{
			name: "admin deletes task -> 200", method: "DELETE", path: "/tasks/t1", token: s.admin,
			setup:  func() { s.tasks.On("DeleteTask", mock.Anything, "t1").Return(nil).Once() },
			status: http.StatusOK,
		},
		{
			name: "missing or someone else's task -> 404", method: "DELETE", path: "/tasks/zzz", token: s.user,
			setup:  func() { s.tasks.On("DeleteTask", mock.Anything, "zzz").Return(Domain.ErrNotFound).Once() },
			status: http.StatusNotFound,
		},
		{
			name: "no token -> 401", method: "DELETE", path: "/tasks/t1",
			status: http.StatusUnauthorized,
		},
	})
}
