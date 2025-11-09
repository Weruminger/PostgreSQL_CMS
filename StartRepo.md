# DB‑First Headless CMS – Startrepo

> PostgreSQL als Kern (RLS/JSONB/Trigger), PostgREST als dünne API, Go‑Portal mit SSR. Auth **wahlweise** anonym/mTLS/Shared‑Secret‑JWT **oder** OIDC/JWT über Keycloak. Minimal‑JS, Multi‑Tenant, Multi‑Domain/Multi‑Portal.

---

## Repo‑Struktur

```
headless-db-cms/
├─ db/
│  ├─ 01_extensions.sql
│  ├─ 02_schema_core.sql
│  ├─ 03_rls_policies.sql
│  ├─ 04_api_views.sql
│  ├─ 05_triggers_functions.sql
│  ├─ 06_media.sql              # NEU: Medienverwaltung
│  ├─ 07_audit.sql              # NEU: Audit-Logs & Diffs
│  └─ 90_demo_seed.sql
├─ postgrest/
│  └─ postgrest.conf
├─ portal/
│  ├─ cmd/portal/main.go
│  └─ internal/
│     ├─ db.go
│     ├─ tenant.go
│     ├─ render.go
│     ├─ auth.go
│     ├─ sessions.go
│     ├─ admin_handlers.go      # NEU
│     ├─ admin_templates.go     # NEU
│     ├─ pgrest.go              # NEU
│     └─ media.go               # NEU: Upload/Proxy
├─ portal/templates/            # NEU: SSR-Templates
│  ├─ admin/
│  │  ├─ layout.tmpl
│  │  ├─ login.tmpl
│  │  ├─ dashboard.tmpl
│  │  ├─ entries_list.tmpl
│  │  ├─ entry_edit.tmpl
│  │  ├─ media_list.tmpl
│  │  ├─ media_upload.tmpl
│  │  └─ audit_list.tmpl
│  └─ public/
│     └─ base.tmpl              # Standard-Theme (Portal)
├─ portal/assets/               # NEU: Minimal-CSS/Icons
│  ├─ admin.css
│  └─ theme.css
├─ docker-compose.yml
├─ .env.example
├─ Makefile
└─ README.md
```

---

## `.env.example`

```dotenv
# Postgres
PGHOST=postgres
PGPORT=5432
PGDATABASE=cms
PGUSER=cms_app
PGPASSWORD=changeme

# App settings (used by portal)
APP_HTTP_ADDR=:8080
APP_CACHE_TTL_SECONDS=30

# Auth mode: "anon" | "sharedjwt" | "keycloak"
AUTH_MODE=anon

# Shared JWT (HS256) - only if AUTH_MODE=sharedjwt
JWT_HS256_SECRET=devdevdev
JWT_AUD=portal
JWT_ISS=self

# Keycloak – only if AUTH_MODE=keycloak
KC_REALM=cms
KC_URL=http://keycloak:8080/realms/cms
KC_AUD=portal

# Admin OIDC Client
OIDC_ISSUER_URL=http://keycloak:8080/realms/cms
OIDC_CLIENT_ID=portal
OIDC_CLIENT_SECRET=dev-secret
OIDC_REDIRECT_URL=https://admin.localhost:8080/admin/callback
SESSION_SECRET=change-me-32bytes-min

# Media storage
MEDIA_BACKEND=fs           # fs | s3 | webdav | cifs
MEDIA_FS_ROOT=/srv/media   # mount in Compose
MEDIA_PUBLIC_BASE=/media   # public URI prefix when proxied through portal
S3_ENDPOINT=https://s3.local
S3_BUCKET=cms-media
S3_ACCESS_KEY=xxx
S3_SECRET_KEY=yyy
WEBDAV_URL=https://webdav.example.com/remote.php/dav/files/user
WEBDAV_USER=user
WEBDAV_PASS=pass
CIFS_MOUNT=//nas/share
CIFS_USER=user
CIFS_PASS=pass
```

---

## `portal/assets/admin.css` (modern, rund, Verläufe)

```css
:root{
  --bg: #0b1220; --panel:#111a2b; --muted:#7a8aa0; --fg:#e8eef7; --acc1:#6ee7f2; --acc2:#a78bfa; --acc3:#34d399;
  --radius: 16px; --radius-sm: 10px; --shadow: 0 8px 30px rgba(0,0,0,.35);
}
*{box-sizing:border-box} body{margin:0; font:16px/1.5 Inter,system-ui,-apple-system,Segoe UI,Roboto; background:
  radial-gradient(1200px 800px at 80% -10%, #1b2540 0%, transparent 60%),
  radial-gradient(900px 600px at -10% 20%, #18223a 0%, transparent 55%),
  linear-gradient(180deg,#0b1220,#0b1220)}
.container{max-width:1200px;margin:0 auto;padding:32px}
.card{background:linear-gradient(180deg,rgba(255,255,255,.03),rgba(255,255,255,.01)); border:1px solid rgba(255,255,255,.08);
  border-radius:var(--radius); box-shadow:var(--shadow); backdrop-filter: blur(8px); color:var(--fg)}
.card h1,.card h2{margin:0 0 12px 0}
.header{display:flex;gap:14px;align-items:center;justify-content:space-between;margin:0 0 24px 0}
.logo{font-weight:700;color:var(--fg)} .muted{color:var(--muted)}
.row{display:grid;gap:16px} .grid-2{grid-template-columns:1fr 1fr} .grid-3{grid-template-columns:repeat(3,1fr)}
.btn{display:inline-flex;align-items:center;gap:8px;padding:10px 14px;border-radius:var(--radius-sm);
  border:1px solid rgba(255,255,255,.12); background:linear-gradient(180deg,rgba(110,231,242,.25),rgba(167,139,250,.25));
  color:var(--fg); text-decoration:none; cursor:pointer}
.btn:hover{filter:brightness(1.05)} .btn.primary{background:linear-gradient(180deg,var(--acc2),var(--acc1))}
.input, select, textarea{width:100%;padding:12px 14px;border-radius:var(--radius-sm);border:1px solid rgba(255,255,255,.12);
  background:rgba(255,255,255,.04); color:var(--fg)}
.table{width:100%;border-collapse:separate;border-spacing:0 8px}
.table tr{background:rgba(255,255,255,.03)}
.table td,.table th{padding:12px 14px}
.badge{padding:4px 10px;border-radius:999px;font-size:12px;border:1px solid rgba(255,255,255,.18)}
.badge.green{background:linear-gradient(180deg,rgba(52,211,153,.25),rgba(52,211,153,.15))}
.badge.yellow{background:linear-gradient(180deg,rgba(234,179,8,.25),rgba(234,179,8,.15))}
.badge.gray{background:linear-gradient(180deg,rgba(148,163,184,.25),rgba(148,163,184,.15))}
.nav{display:flex;gap:8px;flex-wrap:wrap}
.card.pad{padding:20px}
.form-row{display:grid;grid-template-columns:200px 1fr;align-items:center;gap:12px;margin:10px 0}
.code{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace; font-size:13px; background:rgba(0,0,0,.35); padding:8px 10px; border-radius:8px}
```

---

## `portal/assets/theme.css` (Public Theme – Verläufe, runde Ecken)

```css
:root{--radius:18px; --fg:#0e1930; --pri:#2748ff; --sec:#00d4ff}
body{margin:0;font:16px/1.6 Inter,system-ui;-webkit-font-smoothing:antialiased}
.hero{padding:64px 20px;background:linear-gradient(135deg, rgba(39,72,255,.18), rgba(0,212,255,.18));}
.wrap{max-width:1100px;margin:0 auto}
.card{background:#fff;border-radius:var(--radius);box-shadow:0 12px 35px rgba(10,20,50,.08);padding:24px}
.btn{display:inline-block;padding:10px 16px;border-radius:12px;background:linear-gradient(135deg,#2748ff,#00d4ff);color:white;text-decoration:none}
```

---

## Admin‑Templates (SSR)

### `portal/templates/admin/layout.tmpl`

```gotemplate
{{ define "layout" }}
<!doctype html>
<html lang="de">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>{{ block "title" . }}Admin{{ end }}</title>
  <link rel="stylesheet" href="/assets/admin.css" />
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="logo">CMS Admin<span class="muted"> – {{ .TenantSlug }}</span></div>
      <nav class="nav">
        <a class="btn" href="/admin">Dashboard</a>
        <a class="btn" href="/admin/entries">Inhalte</a>
        <a class="btn" href="/admin/media">Medien</a>
        <a class="btn" href="/admin/audit">Audit</a>
        <a class="btn" href="/admin/logout">Logout</a>
      </nav>
    </div>
    {{ block "content" . }}{{ end }}
  </div>
</body>
</html>
{{ end }}
```

### `portal/templates/admin/login.tmpl`

```gotemplate
{{ define "title" }}Login{{ end }}
{{ define "content" }}
<div class="card pad" style="max-width:520px;margin:60px auto">
  <h1>Anmelden</h1>
  <p class="muted">Bitte melde dich über den Identitätsprovider an.</p>
  <a class="btn primary" href="/admin/login/oidc">Mit Keycloak anmelden</a>
</div>
{{ end }}
```

### `portal/templates/admin/dashboard.tmpl`

```gotemplate
{{ define "title" }}Dashboard{{ end }}
{{ define "content" }}
<div class="row grid-3">
  <div class="card pad">
    <h2>Entwürfe</h2>
    <div class="muted">{{ .Stats.Drafts }} offene Entwürfe</div>
  </div>
  <div class="card pad">
    <h2>Zur Prüfung</h2>
    <div class="muted">{{ .Stats.InReview }} Inhalte</div>
  </div>
  <div class="card pad">
    <h2>Veröffentlicht</h2>
    <div class="muted">{{ .Stats.Published }} Artikel</div>
  </div>
</div>
{{ end }}
```

### `portal/templates/admin/entries_list.tmpl`

```gotemplate
{{ define "title" }}Inhalte{{ end }}
{{ define "content" }}
<div class="card pad">
  <div class="header"><h1>Inhalte</h1><a class="btn primary" href="/admin/entries/new">Neu</a></div>
  <table class="table">
    <thead><tr><th>Slug</th><th>Status</th><th>Aktualisiert</th><th></th></tr></thead>
    <tbody>
      {{ range .Items }}
      <tr>
        <td class="code">{{ .Slug }}</td>
        <td>
          {{ if eq .Status "published" }}<span class="badge green">veröffentlicht</span>{{ else if eq .Status "review" }}<span class="badge yellow">zur Prüfung</span>{{ else }}<span class="badge gray">{{ .Status }}</span>{{ end }}
        </td>
        <td>{{ .UpdatedAt }}</td>
        <td><a class="btn" href="/admin/entries/{{ .ID }}">Bearbeiten</a></td>
      </tr>
      {{ end }}
    </tbody>
  </table>
</div>
{{ end }}
```

### `portal/templates/admin/entry_edit.tmpl`

```gotemplate
{{ define "title" }}Eintrag bearbeiten{{ end }}
{{ define "content" }}
<form method="post" action="/admin/entries/{{ .Entry.ID }}" class="card pad">
  <input type="hidden" name="csrf" value="{{ .CSRF }}" />
  <div class="form-row"><label>Titel</label><input class="input" name="title" value="{{ .Entry.Title }}" /></div>
  <div class="form-row"><label>Body</label><textarea class="input" name="body" rows="10">{{ .Entry.Body }}</textarea></div>
  <div class="form-row"><label>Status</label>
    <select name="status">
      <option value="draft" {{ if eq .Entry.Status "draft" }}selected{{ end }}>Entwurf</option>
      <option value="review" {{ if eq .Entry.Status "review" }}selected{{ end }}>Zur Prüfung</option>
      <option value="published" {{ if eq .Entry.Status "published" }}selected{{ end }}>Veröffentlicht</option>
    </select>
  </div>
  <div class="form-row"><label>Medien</label>
    <div>
      {{ range .LinkedMedia }}<span class="badge">{{ .Filename }}</span>{{ end }}
      <a class="btn" href="/admin/media?link={{ .Entry.ID }}">Medien verknüpfen</a>
    </div>
  </div>
  <div style="display:flex; gap:10px; justify-content:flex-end">
    <button class="btn" name="action" value="save" type="submit">Speichern</button>
    <button class="btn primary" name="action" value="publish" type="submit">Freigeben</button>
  </div>
</form>
{{ end }}
```

### `portal/templates/admin/media_list.tmpl`

```gotemplate
{{ define "title" }}Medien{{ end }}
{{ define "content" }}
<div class="card pad">
  <div class="header"><h1>Medien</h1><a class="btn" href="/admin/media/upload">Upload</a></div>
  <table class="table"><thead><tr><th>Datei</th><th>MIME</th><th>Größe</th><th>Checksumme</th><th></th></tr></thead>
  <tbody>
    {{ range .Items }}
    <tr>
      <td>{{ .Filename }}</td>
      <td>{{ .Mime }}</td>
      <td>{{ .SizeHuman }}</td>
      <td class="code">{{ .Checksum }}</td>
      <td>
        <a class="btn" href="/media/{{ .StorageKey }}" target="_blank">Ansehen</a>
        {{ if $.LinkTo }}<a class="btn" href="/admin/media/link?asset={{ .ID }}&entry={{ $.LinkTo }}">Verknüpfen</a>{{ end }}
      </td>
    </tr>
    {{ end }}
  </tbody></table>
</div>
{{ end }}
```

### `portal/templates/admin/media_upload.tmpl`

```gotemplate
{{ define "title" }}Upload{{ end }}
{{ define "content" }}
<form class="card pad" method="post" action="/admin/media/upload" enctype="multipart/form-data">
  <input type="hidden" name="csrf" value="{{ .CSRF }}" />
  <div class="form-row"><label>Datei</label><input class="input" type="file" name="file" /></div>
  <div class="form-row"><label>Alt‑Text</label><input class="input" name="alt" /></div>
  <div style="display:flex; gap:10px; justify-content:flex-end"><button class="btn primary" type="submit">Hochladen</button></div>
</form>
{{ end }}
```

### `portal/templates/admin/audit_list.tmpl`

```gotemplate
{{ define "title" }}Audit Log{{ end }}
{{ define "content" }}
<div class="card pad">
  <h1>Änderungen</h1>
  <table class="table"><thead><tr><th>Zeit</th><th>User</th><th>Entität</th><th>ID</th><th>Aktion</th><th>Diff</th></tr></thead>
    <tbody>
    {{ range .Items }}
      <tr>
        <td>{{ .At }}</td><td>{{ .Actor }}</td><td>{{ .Entity }}</td><td>{{ .EntityID }}</td><td>{{ .Action }}</td>
        <td><pre class="code">{{ .Diff }}</pre></td>
      </tr>
    {{ end }}
    </tbody>
  </table>
</div>
{{ end }}
```

---

## Public Theme (Portal) – `portal/templates/public/base.tmpl`

```gotemplate
{{ define "layout" }}
<!doctype html>
<html lang="de"><head>
<meta charset="utf-8" /><meta name="viewport" content="width=device-width, initial-scale=1" />
<link rel="stylesheet" href="/assets/theme.css" />
<title>{{ .data.title }}</title></head>
<body>
  <section class="hero"><div class="wrap"><h1>{{ .data.title }}</h1></div></section>
  <main class="wrap"><div class="card">{{ .data.body }}</div></main>
</body></html>
{{ end }}
```

---

## Medien‑Schema & Backend

### `db/06_media.sql`

```sql
create table if not exists media_locations (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id) on delete cascade,
  backend text not null check (backend in ('fs','s3','webdav','cifs')),
  base_uri text not null,         -- z.B. file:///srv/media, s3://bucket/prefix, https://webdav/…
  config jsonb not null default '{}'::jsonb,
  is_default boolean not null default false
);

create table if not exists media_assets (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id) on delete cascade,
  location_id bigint not null references media_locations(id) on delete restrict,
  filename text not null,
  mime_type text not null,
  size_bytes bigint not null,
  checksum text not null,         -- sha256
  storage_key text not null,      -- Pfad/Key relativ zur base_uri
  alt text,
  meta jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);

create table if not exists entry_media (
  entry_id bigint not null references entries(id) on delete cascade,
  asset_id bigint not null references media_assets(id) on delete cascade,
  primary key (entry_id, asset_id)
);

-- RLS
alter table media_locations enable row level security;
alter table media_assets enable row level security;
alter table entry_media enable row level security;

create policy media_locations_tenant on media_locations using (tenant_id = app_current_tenant_id()) with check (tenant_id = app_current_tenant_id());
create policy media_assets_tenant   on media_assets   using (tenant_id = app_current_tenant_id()) with check (tenant_id = app_current_tenant_id());
create policy entry_media_tenant    on entry_media    using (
  exists (select 1 from entries e where e.id = entry_id and e.tenant_id = app_current_tenant_id())
) with check (
  exists (select 1 from entries e where e.id = entry_id and e.tenant_id = app_current_tenant_id())
);

-- Grants
grant select on media_assets, media_locations to web_anon;
```

### Portal‑Upload/Proxy (Skizze) – `portal/internal/media.go`

```go
// Pseudocode: je nach MEDIA_BACKEND speichern
// fs: in MEDIA_FS_ROOT/<tenant>/<uuid> ablegen
// s3: PutObject(S3_BUCKET, key)
// webdav/cifs: Pfad ist gemountet, wir schreiben wie fs
```

---

## Audit‑Logs & Diffs

### `db/07_audit.sql`

```sql
create table if not exists admin_audit (
  id bigserial primary key,
  tenant_id bigint not null,
  actor text not null,              -- aus Token (sub/email)
  entity text not null,             -- 'entries' | 'themes' | 'media_assets' ...
  entity_id bigint not null,
  action text not null,             -- 'insert' | 'update' | 'delete'
  diff jsonb,                       -- JSON Patch/Delta (optional)
  at timestamptz not null default now()
);

alter table admin_audit enable row level security;
create policy audit_tenant on admin_audit for select using (tenant_id = app_current_tenant_id());

-- Hilfsfunktion: einfache Feld-Diff für JSONB (ohne plpython)
create or replace function jsonb_shallow_diff(old jsonb, new jsonb)
returns jsonb language sql immutable as $$
  select coalesce(jsonb_object_agg(k, v), '{}'::jsonb) from (
    select k, new->k as v
    from (select jsonb_object_keys(coalesce(old,'{}'::jsonb) || coalesce(new,'{}'::jsonb)) as k) ks
    where coalesce(old->k,'null'::jsonb) is distinct from coalesce(new->k,'null'::jsonb)
  ) s;
$$;

-- Trigger für entries (schreibt Audit + Versionen gibt es bereits)
create or replace function audit_entries() returns trigger language plpgsql as $$
begin
  if TG_OP='INSERT' then
    insert into admin_audit(tenant_id, actor, entity, entity_id, action, diff)
    values (new.tenant_id, current_setting('request.jwt.claims',true)::jsonb->>'sub', 'entries', new.id, 'insert', new.data);
  elsif TG_OP='UPDATE' then
    insert into admin_audit(tenant_id, actor, entity, entity_id, action, diff)
    values (new.tenant_id, current_setting('request.jwt.claims',true)::jsonb->>'sub', 'entries', new.id, 'update', jsonb_shallow_diff(old.data, new.data));
  elsif TG_OP='DELETE' then
    insert into admin_audit(tenant_id, actor, entity, entity_id, action, diff)
    values (old.tenant_id, current_setting('request.jwt.claims',true)::jsonb->>'sub', 'entries', old.id, 'delete', old.data);
  end if;
  return COALESCE(new, old);
end; $$;

create trigger trg_audit_entries
after insert or update or delete on entries
for each row execute procedure audit_entries();

-- Optional (präziser): JSON Patch mit plpython3u + jsonpatch, wenn erlaubt
-- create function jsonb_patch(old jsonb, new jsonb) returns jsonb language plpython3u ...
```

### Admin‑Ansicht

* `/admin/audit`: listet `admin_audit` (RLS filtert auf Mandant). Diff wird als JSON angezeigt; für „schön“ kann serverseitig ein kleines Pretty‑Printer‑Template genutzt werden.

---

## Admin‑Handlers – Routen (Skizze)

* **GET** `/admin/media` → Liste aus `media_assets` (Tenant via RLS), optional `?link=<entry_id>`
* **GET/POST** `/admin/media/upload` → multipart Upload → Datei per Backend‑Treiber speichern, anschließend `media_assets`‑Zeile anlegen.
* **GET** `/admin/entries` → Liste
* **GET** `/admin/entries/{id}` → Formular
* **POST** `/admin/entries/{id}` → `PATCH /entries` via PostgREST mit Bearer‑Token
* **POST** `/admin/media/link?asset=X&entry=Y` → `insert into entry_media`

Sämtliche Schreibzugriffe gehen **nur** über PostgREST (Bearer‑Token aus Session), wodurch RLS/Policies/Audit‑Trigger greifen.

---

## Hinweise zur Produktion

* **Assets liefern**: Medien bevorzugt über ein separates CDN/Reverse‑Proxy. Für NFS/CIFS Mounds im Container read‑write nur im Admin, Public‑Portal nur read‑only/Proxy.
* **S3/WebDAV**: Credentials niemals im Image, nur über Secrets/Env. Signierte URLs optional (Zeitfenster zum Download).
* **Audit‑Tiefe**: Für vollständige JSON‑Patch‑Deltas plpython3u + `jsonpatch` verwenden; ansonsten `jsonb_shallow_diff` reicht für viele CMS‑Fälle (Titel/Teaser/Body etc.).
* **Backpressure** beim Upload: Max‑Größe limitieren, Antivirus/ClamAV‑Hook optional.
