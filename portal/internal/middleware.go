package internal

import (
	"context"
	"net/http"
)

type ctxKey string

const dbKey ctxKey = "db"

func WithDB(db *DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), dbKey, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getDB(r *http.Request) *DB {
	if v := r.Context().Value(dbKey); v != nil {
		if db, ok := v.(*DB); ok {
			return db
		}
	}
	return nil
}
