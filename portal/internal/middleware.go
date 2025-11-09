package internal
import (
	"context"
	"net/http"
)
type ctxKey string
const dbKey ctxKey = "db"
func WithDB(db *DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), dbKey, db)))
	})
}
func getDB(r *http.Request) *DB { v := r.Context().Value(dbKey); if v==nil { return nil }; return v.(*DB) }
func contextWithDB(r *http.Request) *http.Request { return r }
