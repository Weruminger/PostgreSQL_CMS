package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("/srv/templates/admin/layout.tmpl", "/srv/templates/admin/login.tmpl"))
	_ = tpl.ExecuteTemplate(w, "layout", map[string]any{"TenantSlug": "demo"})
}

func AdminLoginOIDC(w http.ResponseWriter, r *http.Request) {
	issuer := getenv("OIDC_ISSUER_URL", "")
	clientID := getenv("OIDC_CLIENT_ID", "portal")
	redir := getenv("OIDC_REDIRECT_URL", "http://localhost:8080/admin/callback")
	authURL := fmt.Sprintf("%s/protocol/openid-connect/auth?response_type=code&client_id=%s&redirect_uri=%s&scope=openid",
		issuer, url.QueryEscape(clientID), url.QueryEscape(redir))
	http.Redirect(w, r, authURL, http.StatusFound)
}

func AdminCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}
	issuer := getenv("OIDC_ISSUER_URL", "")
	tokenURL := issuer + "/protocol/openid-connect/token"
	redir := getenv("OIDC_REDIRECT_URL", "http://localhost:8080/admin/callback")
	clientID := getenv("OIDC_CLIENT_ID", "portal")
	clientSecret := getenv("OIDC_CLIENT_SECRET", "dev-secret")

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redir)
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		http.Error(w, "oidc", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	var tok map[string]any
	_ = json.Unmarshal(b, &tok)
	at, _ := tok["access_token"].(string)

	role := "editor"
	tenantID := int64(1)

	if at != "" {
		if pr, err := parseJWTPayload(at); err == nil {
			if ra, ok := pr["realm_access"].(map[string]any); ok {
				if roles, ok := ra["roles"].([]any); ok {
					for _, r := range roles {
						if rs, ok := r.(string); ok && (rs == "admin" || rs == "editor") {
							role = rs
							break
						}
					}
				}
			}
		}
	}

	setSession(w, r, map[string]any{
		"access_token": at,
		"pgrst_role":   role,
		"tenant_id":    tenantID,
	})
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := getSession(r); !ok {
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func parseJWTPayload(token string) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("bad token")
	}
	p, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var m map[string]any
	return m, json.Unmarshal(p, &m)
}
