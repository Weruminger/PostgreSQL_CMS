package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithDBAndGetDB(t *testing.T) {
	db := &DB{}
	h := WithDB(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := getDB(r)
		require.Same(t, db, got)
	}))
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
}
