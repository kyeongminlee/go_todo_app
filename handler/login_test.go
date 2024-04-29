package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-playground/validator"
	"go_todo_app/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin_ServedHTTP(t *testing.T) {
	type moq struct {
		token string
		err   error
	}
	type want struct {
		status       int
		responseFile string
	}
	tests := map[string]struct {
		requestFile string
		moq         moq
		want        want
	}{
		"ok": {
			requestFile: "testdata/login/ok_request.json.golden",
			moq: moq{
				token: "from_moq",
			},
			want: want{
				status:       http.StatusOK,
				responseFile: "testdata/login/ok_response.json.golden",
			},
		},
		"badRequest": {
			requestFile: "testdata/login/bad_request.json.golden",
			want: want{
				status:       http.StatusBadRequest,
				responseFile: "testdata/login/bad_response.json.golden",
			},
		},
		"internal_server_error": {
			requestFile: "testdata/login/ok_request.json.golden",
			moq: moq{
				err: errors.New("error from mock"),
			},
			want: want{
				status:       http.StatusInternalServerError,
				responseFile: "testdata/login/internal_server_error_response.json.golden",
			},
		},
	}
	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewReader(testutil.LoadFile(t, tt.requestFile)),
			)

			moq := &LoginServiceMock{}
			moq.LoginFunc = func(ctx context.Context, name, pw string) (string, error) {
				return tt.moq.token, tt.moq.err
			}
			sut := Login{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServedHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.responseFile))
		})
	}
}
