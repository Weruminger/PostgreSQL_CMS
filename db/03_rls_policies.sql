create or replace function app_current_tenant_id() returns bigint language sql stable as $$
  select coalesce(
    nullif(current_setting('app.tenant_id', true), '')::bigint,
    (current_setting('request.jwt.claims', true)::jsonb ->> 'tenant_id')::bigint
  );
$$;
create or replace function app_current_role() returns text language sql stable as $$
  select coalesce(
    nullif(current_setting('app.role', true), ''),
    (current_setting('request.jwt.claims', true)::jsonb ->> 'pgrst_role'),
    'public'
  );
$$;
alter table tenants enable row level security;
alter table domains enable row level security;
alter table themes enable row level security;
alter table content_types enable row level security;
alter table entries enable row level security;
alter table entry_versions enable row level security;
create policy tenants_iso on tenants
  for all using (id = app_current_tenant_id()) with check (id = app_current_tenant_id());
create policy domains_iso on domains
  for all using (tenant_id = app_current_tenant_id()) with check (tenant_id = app_current_tenant_id());
create policy themes_iso on themes
  for all using (tenant_id = app_current_tenant_id()) with check (tenant_id = app_current_tenant_id());
create policy ctypes_iso on content_types
  for all using (tenant_id = app_current_tenant_id()) with check (tenant_id = app_current_tenant_id());
create policy entries_iso on entries
  for all using (tenant_id = app_current_tenant_id()) with check (tenant_id = app_current_tenant_id());
create policy entries_public_read on entries
  for select using (app_current_role() in ('public','web_anon') and status='published');
create policy entries_edit_write on entries
  for all using (app_current_role() in ('editor','admin'))
  with check (app_current_role() in ('editor','admin'));
grant usage on schema public to web_anon;
grant select on entries, domains, themes, content_types, tenants to web_anon;
grant select, insert, update, delete on all tables in schema public to cms_app;
alter default privileges in schema public grant select, insert, update, delete on tables to cms_app;
