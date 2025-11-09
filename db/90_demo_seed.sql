insert into tenants (slug, title) values ('demo', 'Demo Tenant') returning id \gset
insert into domains (tenant_id, host, primary_domain) values (:id, 'localhost', true);
insert into themes (tenant_id, name, layout, partials) values
(:id, 'default',
$${{ define "layout" }}
<!doctype html><html><head><meta charset="utf-8"><title>{{ .data.title }}</title>
<link rel="stylesheet" href="/assets/theme.css"></head>
<body><section class="hero"><div class="wrap"><h1>{{ .data.title }}</h1></div></section>
<main class="wrap"><div class="card">{{ .data.body }}</div></main></body></html>
{{ end }}$$,
json_build_object('header','<header><a href="/">Demo</a></header>')
);
insert into content_types (tenant_id, slug, schema) values
(:id, 'article', '{"type":"object","properties":{"title":{"type":"string"},"body":{"type":"string"}},"required":["title","body"]}');
insert into entries (tenant_id, type_id, slug, data, status, published_at)
select :id, ct.id, 'hello-world', '{"title":"Hello","body":"World"}', 'published', now()
from content_types ct where ct.slug='article' and ct.tenant_id=:id;
insert into media_locations(tenant_id, backend, base_uri, is_default) values (:id, 'fs', 'file:///srv/media', true);
