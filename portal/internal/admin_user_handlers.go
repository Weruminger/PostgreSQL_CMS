package internal

import (
  "context"
  "encoding/json"
  "net/http"
  "strconv"
)

// AdminUserDelete handles tenant-admin initiated user deletion (GDPR)
// Route example: POST /admin/tenants/{tenant_id}/users/{user_id}/delete
func AdminUserDelete(w http.ResponseWriter, r *http.Request) {
  db := getDB(r)

  // naive path parsing (adapt to your router)
  // expected path: /admin/tenants/{tenant_id}/users/{user_id}/delete
  parts := splitPath(r.URL.Path)
  if len(parts) < 6 {
    http.Error(w, "invalid path", 400)
    return
  }
  tenantID, err := strconv.ParseInt(parts[2], 10, 64)
  if err != nil {
    http.Error(w, "invalid tenant id", 400)
    return
  }
  userID, err := strconv.ParseInt(parts[4], 10, 64)
  if err != nil {
    http.Error(w, "invalid user id", 400)
    return
  }

  var report map[string]any
  err = db.WithAppSettings(r.Context(), tenantID, "admin", func(ctx context.Context, q Querier) error {
    var res string
    err := q.QueryRow(ctx, `select delete_user_gdpr($1,$2)::text`, tenantID, userID).Scan(&res)
    if err != nil { return err }
    return json.Unmarshal([]byte(res), &report)
  })
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  enc := json.NewEncoder(w)
  enc.Encode(report)
}

// splitPath is a small helper to split URL path cleanly
func splitPath(p string) []string {
  // remove leading and trailing /
  for len(p) > 0 && p[0] == '/' { p = p[1:] }
  for len(p) > 0 && p[len(p)-1] == '/' { p = p[:len(p)-1] }
  if p == "" { return []string{} }
  return split(p, '/')
}

// simple split to avoid importing strings in this stub
func split(s string, sep rune) []string {
  var res []string
  cur := ""
  for _, r := range s {
    if r == sep { res = append(res, cur); cur = ""; continue }
    cur += string(r)
  }
  res = append(res, cur)
  return res
}