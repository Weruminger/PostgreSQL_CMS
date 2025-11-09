package internal

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionSetGetClear(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	setSession(w, r, map[string]any{"u": "abc"})

	resp := w.Result()
	cookies := resp.Cookies()
	require.NotEmpty(t, cookies)
	c := cookies[0]

	r2 := httptest.NewRequest("GET", "/", nil)
	r2.AddCookie(c)
	m, ok := getSession(r2)
	require.True(t, ok)
	require.Equal(t, "abc", m["u"])

	w2 := httptest.NewRecorder()
	clearSession(w2)
	require.NotEmpty(t, w2.Result().Cookies())
}
