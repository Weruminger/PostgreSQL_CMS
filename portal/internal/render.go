package internal
import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)
type theme struct { Layout string; Partials map[string]string }
type cacheEntry struct { th theme; exp time.Time }
var themeCache = map[int64]cacheEntry{}
func PageHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		path := strings.Trim(r.URL.Path, "/")
		if strings.HasPrefix(path, "admin") || strings.HasPrefix(path, "assets") { http.NotFound(w,r); return }
		ctx := r.Context()
		tenantID := int64(1)
		if id, err := LookupTenantByHost(ctx, db.Pool(), host); err==nil { tenantID = id }
		err := db.WithAppSettings(ctx, tenantID, "public", func(ctx context.Context, q Querier) error {
			if path == "" { path = "hello-world" }
			art, err := LoadArticleBySlug(ctx, q, path); if err!=nil { http.NotFound(w,r); return nil }
			th, err := getTheme(ctx, q, tenantID); if err!=nil { http.Error(w,"theme",500); return nil }
			tpl := template.Must(template.New("layout").Parse(th.Layout))
			for name, part := range th.Partials { tpl = template.Must(tpl.New(name).Parse(part)) }
			w.Header().Set("Content-Type","text/html; charset=utf-8")
			return tpl.ExecuteTemplate(w, "layout", map[string]any{"data": map[string]any{"title": art.Title, "body": art.Body}})
		}); if err!=nil { log.Printf("page: %v", err); http.Error(w,"server",500) }
	}
}
func getTheme(ctx context.Context, q Querier, tenantID int64) (theme, error) {
	if ce, ok := themeCache[tenantID]; ok && time.Now().Before(ce.exp) { return ce.th, nil }
	var layout, partialsJSON string
	err := q.QueryRow(ctx, `select layout, partials::text from themes where tenant_id=$1 order by id desc limit 1`, tenantID).Scan(&layout, &partialsJSON)
	if err != nil { return theme{}, err }
	partials := map[string]string{}
	_ = json.Unmarshal([]byte(partialsJSON), &partials)
	th := theme{Layout:layout, Partials:partials}
	themeCache[tenantID] = cacheEntry{ th: th, exp: time.Now().Add(30*time.Second) }
	return th, nil
}
