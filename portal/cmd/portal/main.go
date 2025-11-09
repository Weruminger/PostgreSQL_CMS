package main
import (
	"log"
	"net/http"
	"os"
	"time"
	"headless-db-cms/internal"
)
func main(){
	addr := getenv("APP_HTTP_ADDR", ":8080")
	db, err := internal.NewDBFromEnv()
	if err != nil { log.Fatalf("db: %v", err) }
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/", internal.PageHandler(db))
	mux.HandleFunc("/admin/login", internal.AdminLoginPage)
	mux.HandleFunc("/admin/login/oidc", internal.AdminLoginOIDC)
	mux.HandleFunc("/admin/callback", internal.AdminCallback)
	mux.HandleFunc("/admin/logout", internal.AdminLogout)
	mux.Handle("/admin", internal.RequireAuth(http.HandlerFunc(internal.AdminDashboard)))
	mux.Handle("/admin/entries", internal.RequireAuth(http.HandlerFunc(internal.AdminEntriesList)))
	mux.Handle("/admin/entries/new", internal.RequireAuth(http.HandlerFunc(internal.AdminEntryNew)))
	mux.Handle("/admin/entries/", internal.RequireAuth(http.HandlerFunc(internal.AdminEntryEdit)))
	mux.Handle("/admin/media", internal.RequireAuth(http.HandlerFunc(internal.AdminMediaList)))
	mux.Handle("/admin/media/upload", internal.RequireAuth(http.HandlerFunc(internal.AdminMediaUpload)))
	mux.Handle("/admin/audit", internal.RequireAuth(http.HandlerFunc(internal.AdminAuditList)))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("/srv/assets"))))
	srv := &http.Server{ Addr: addr, Handler: internal.WithDB(db, internal.AuthMiddleware(mux)), ReadTimeout: 5*time.Second, WriteTimeout: 15*time.Second }
	log.Printf("portal listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
func getenv(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }
