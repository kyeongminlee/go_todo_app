package handler

import (
	"bytes"
	"context"
	"errors"
	"go_todo_app/entity"
	"go_todo_app/testutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
)

func TestAddTask(t *testing.T) {
	type want struct {
		status       int
		responseFile string
	}
	tests := map[string]struct {
		requestFile string
		want        want
	}{
		"ok": {
			requestFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:       http.StatusOK,
				responseFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			requestFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:       http.StatusBadRequest,
				responseFile: "testdata/add_task/bad_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			writer := httptest.NewRecorder()
			request := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.requestFile)),
			)

			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				if tt.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			sut := AddTask{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(writer, request)

			response := writer.Result()
			testutil.AssertResponse(t,
				response, tt.want.status, testutil.LoadFile(t, tt.want.responseFile))

		})
	}
}
