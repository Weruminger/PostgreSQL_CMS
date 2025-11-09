package internal

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionSetGetClear(t *testing.T) {
	// set
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	setSession(w, r, map[string]any{"u": "abc"})

	// simulate browser cookie roundtrip
	resp := w.Result()
	c := resp.Cookies()[0]

	// get
	r2 := httptest.NewRequest("GET", "/", nil)
	r2.AddCookie(c)
	m, ok := getSession(r2)
	require.True(t, ok)
	require.Equal(t, "abc", m["u"])

	// clear
	w2 := httptest.NewRecorder()
	clearSession(w2)
	require.Contains(t, w2.Result().Cookies()[0].Value, "")
}
