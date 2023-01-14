package weblogin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCookieValue(t *testing.T) {
	cases := []struct {
		request    bool
		cookie     *http.Cookie
		name, want string
		err        error
	}{
		{
			request: false, cookie: nil,
			name: "", want: "", err: ErrNoRequest,
		},
		{
			request: true, cookie: nil,
			name: "none", want: "", err: nil,
		},
		{
			request: true,
			cookie:  &http.Cookie{Name: "test", Value: "value"},
			name:    "test", want: "value", err: nil,
		},
		{
			request: true,
			cookie:  &http.Cookie{Name: "test", Value: "value"},
			name:    "none", want: "", err: nil,
		},
	}

	for _, tc := range cases {
		var r *http.Request

		if tc.request {
			r = httptest.NewRequest(http.MethodGet, "/test", nil)
		}
		if tc.cookie != nil {
			r.AddCookie(tc.cookie)
		}

		got, err := GetCookieValue(r, tc.name)
		if !errors.Is(err, tc.err) {
			t.Errorf("GetCookieValue(%v, %q)\ngot err '%v' want '%v'", r, tc.name, err, tc.err)
		}
		if got != tc.want {
			t.Errorf("GetCookieValue(%v, %q)\ngot %q want %q", r, tc.name, got, tc.want)
		}
	}
}
