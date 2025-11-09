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