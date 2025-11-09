package internal
import (
	"context"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)
func AdminDashboard(w http.ResponseWriter, r *http.Request){
	tpl := template.Must(template.ParseFiles("/srv/templates/admin/layout.tmpl", "/srv/templates/admin/dashboard.tmpl"))
	stats := struct{Drafts,InReview,Published int}{0,0,0}
	tpl.ExecuteTemplate(w,"layout", map[string]any{"TenantSlug":"demo","Stats":stats})
}
func AdminEntriesList(w http.ResponseWriter, r *http.Request){
	tpl := template.Must(template.ParseFiles("/srv/templates/admin/layout.tmpl", "/srv/templates/admin/entries_list.tmpl"))
	db := getDB(r); var items []map[string]any
	db.WithAppSettings(r.Context(), 1, "editor", func(ctx context.Context, q Querier) error{
		rows, _ := q.Query(ctx, `select id, slug, status, to_char(updated_at,'YYYY-MM-DD HH24:MI') from entries where tenant_id=1 order by updated_at desc limit 100`)
		for rows.Next(){ var id int64; var slug,status,updated string; rows.Scan(&id,&slug,&status,&updated); items = append(items, map[string]any{"ID":id,"Slug":slug,"Status":status,"UpdatedAt":updated}) }
		return nil
	})
	tpl.ExecuteTemplate(w,"layout", map[string]any{"TenantSlug":"demo","Items":items})
}
func AdminEntryNew(w http.ResponseWriter, r *http.Request){
	db := getDB(r)
	db.WithAppSettings(r.Context(), 1, "editor", func(ctx context.Context, q Querier) error{
		var typeID int64; q.QueryRow(ctx, `select id from content_types where slug='article' and tenant_id=1`).Scan(&typeID)
		slug := fmt.Sprintf("new-%d", os.Getpid())
		q.Exec(ctx, `insert into entries(tenant_id,type_id,slug,data,status) values(1,$1,$2,'{"title":"Neu","body":""}','draft')`, typeID, slug)
		return nil
	})
	http.Redirect(w,r,"/admin/entries",302)
}
func AdminEntryEdit(w http.ResponseWriter, r *http.Request){
	db := getDB(r); idStr := strings.TrimPrefix(r.URL.Path, "/admin/entries/"); id, _ := strconv.ParseInt(idStr, 10, 64)
	if r.Method=="POST" { title := r.FormValue("title"); body := r.FormValue("body"); status := r.FormValue("status"); db.WithAppSettings(r.Context(), 1, "editor", func(ctx context.Context, q Querier) error{ q.Exec(ctx, `update entries set data=jsonb_set(jsonb_set(data,'{title}',to_jsonb($1::text)), '{body}', to_jsonb($2::text)), status=$3 where id=$4`, title, body, status, id); return nil }); http.Redirect(w,r,"/admin/entries/"+idStr,302); return }
	var title, body, status string
	db.WithAppSettings(r.Context(), 1, "editor", func(ctx context.Context, q Querier) error{ q.QueryRow(ctx, `select data->>'title', data->>'body', status from entries where id=$1`, id).Scan(&title, &body, &status); return nil })
	tpl := template.Must(template.ParseFiles("/srv/templates/admin/layout.tmpl", "/srv/templates/admin/entry_edit.tmpl"))
	tpl.ExecuteTemplate(w, "layout", map[string]any{"TenantSlug":"demo","Entry":map[string]any{"ID":id,"Title":title,"Body":body,"Status":status}, "CSRF":"dev"})
}
func AdminMediaList(w http.ResponseWriter, r *http.Request){
	db := getDB(r); var items []map[string]any
	db.WithAppSettings(r.Context(), 1, "editor", func(ctx context.Context, q Querier) error{
		rows, _ := q.Query(ctx, `select id, filename, mime_type, size_bytes, checksum, storage_key from media_assets where tenant_id=1 order by created_at desc limit 100`)
		for rows.Next(){ var id int64; var fn, mm, cs, key string; var sz int64; rows.Scan(&id,&fn,&mm,&sz,&cs,&key); items = append(items, map[string]any{"ID":id,"Filename":fn,"Mime":mm,"Checksum":cs,"SizeHuman":fmt.Sprintf("%d B",sz),"StorageKey":key}) }
		return nil
	})
	tpl := template.Must(template.ParseFiles("/srv/templates/admin/layout.tmpl", "/srv/templates/admin/media_list.tmpl"))
	tpl.ExecuteTemplate(w, "layout", map[string]any{"TenantSlug":"demo","Items":items})
}
func AdminMediaUpload(w http.ResponseWriter, r *http.Request){
	if r.Method=="GET" { tpl := template.Must(template.ParseFiles("/srv/templates/admin/layout.tmpl", "/srv/templates/admin/media_upload.tmpl")); tpl.ExecuteTemplate(w, "layout", map[string]any{"TenantSlug":"demo","CSRF":"dev"}); return }
	file, header, err := r.FormFile("file"); if err!=nil { http.Error(w,"upload",400); return }
	defer file.Close()
	b, _ := io.ReadAll(file); sum := sha256.Sum256(b)
	name := header.Filename; key := fmt.Sprintf("%s/%s", "1", name) // tenant 1 demo
	root := getenv("MEDIA_FS_ROOT","/srv/media"); path := filepath.Join(root, key)
	os.MkdirAll(filepath.Dir(path), 0o755); os.WriteFile(path, b, 0o644)
	db := getDB(r)
	db.WithAppSettings(r.Context(), 1, "editor", func(ctx context.Context, q Querier) error{
		var locID int64; q.QueryRow(ctx, `select id from media_locations where tenant_id=1 and is_default true limit 1`).Scan(&locID)
		q.Exec(ctx, `insert into media_assets(tenant_id,location_id,filename,mime_type,size_bytes,checksum,storage_key) values(1,$1,$2,$3,$4,$5,$6)`, locID, name, header.Header.Get("Content-Type"), len(b), fmt.Sprintf("%x", sum), key)
		return nil
	})
	http.Redirect(w,r,"/admin/media",302)
}
func AdminAuditList(w http.ResponseWriter, r *http.Request){
	db := getDB(r); var items []map[string]any
	db.WithAppSettings(r.Context(), 1, "editor", func(ctx context.Context, q Querier) error{
		rows, _ := q.Query(ctx, `select at, coalesce(actor,'?'), entity, entity_id, action, to_json(diff)::text from admin_audit where tenant_id=1 order by at desc limit 100`)
		for rows.Next(){ var at, actor, entity, action, diff string; var id int64; rows.Scan(&at,&actor,&entity,&id,&action,&diff); items = append(items, map[string]any{"At":at,"Actor":actor,"Entity":entity,"EntityID":id,"Action":action,"Diff":diff}) }
		return nil
	})
	tpl := template.Must(template.ParseFiles("/srv/templates/admin/layout.tmpl", "/srv/templates/admin/audit_list.tmpl"))
	tpl.ExecuteTemplate(w, "layout", map[string]any{"TenantSlug":"demo","Items":items})
}
