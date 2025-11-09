do $$ begin
  if not exists (select from pg_roles where rolname='cms_app') then
    create role cms_app login password 'changeme';
  end if;
end $$;
create role web_anon  nologin;
create role editor    nologin;
create role admin     nologin;

create table if not exists tenants (
  id bigserial primary key,
  slug text unique not null,
  title text not null,
  active boolean not null default true
);
create table if not exists domains (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id) on delete cascade,
  host text not null,
  primary_domain boolean not null default false,
  unique(tenant_id, host)
);
create table if not exists themes (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id) on delete cascade,
  name text not null,
  layout text not null,
  partials jsonb not null default '{}'::jsonb,
  assets jsonb not null default '{}'::jsonb,
  version int not null default 1
);
create table if not exists content_types (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id) on delete cascade,
  slug text not null,
  schema jsonb not null,
  unique(tenant_id, slug)
);
create table if not exists entries (
  id bigserial primary key,
  tenant_id bigint not null references tenants(id) on delete cascade,
  type_id bigint not null references content_types(id) on delete restrict,
  slug text not null,
  data jsonb not null,
  status text not null check (status in ('draft','review','published','archived')),
  published_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique(tenant_id, type_id, slug)
);
create table if not exists entry_versions (
  id bigserial primary key,
  entry_id bigint not null references entries(id) on delete cascade,
  version int not null,
  data jsonb not null,
  editor text,
  created_at timestamptz not null default now(),
  unique(entry_id, version)
);
alter table entries add column if not exists fts tsvector generated always as
(to_tsvector('simple',coalesce(data->>'title','')||' '||coalesce(data->>'summary','')||' '||coalesce(data->>'body',''))) stored;
create index if not exists entries_fts_idx on entries using gin (fts);
create index if not exists entries_pub_idx on entries (tenant_id, type_id, status, published_at desc);
