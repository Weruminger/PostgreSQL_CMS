package internal
import (
	"encoding/gob"
	"net/http"
	"os"
	"github.com/gorilla/securecookie"
)
var sc = securecookie.New([]byte(getenv("SESSION_SECRET","dev-secret-please-change")), nil)
func init(){ gob.Register(map[string]any{}) }
func setSession(w http.ResponseWriter, r *http.Request, data map[string]any){ val, _ := sc.Encode("sess", data); http.SetCookie(w, &http.Cookie{Name:"sess", Value:val, Path:"/", HttpOnly:true, Secure:false, SameSite:http.SameSiteLaxMode}) }
func getSession(r *http.Request) (map[string]any, bool){ c, err := r.Cookie("sess"); if err!=nil { return nil,false }; var m map[string]any; if err := sc.Decode("sess", c.Value, &m); err!=nil { return nil,false }; return m, true }
func clearSession(w http.ResponseWriter){ http.SetCookie(w, &http.Cookie{Name:"sess", Value:"", Path:"/", MaxAge:-1}) }
func getenv(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }
