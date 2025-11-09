create table if not exists admin_audit (
  id bigserial primary key,
  tenant_id bigint not null,
  actor text not null,
  entity text not null,
  entity_id bigint not null,
  action text not null,
  diff jsonb,
  at timestamptz not null default now()
);
alter table admin_audit enable row level security;
create policy audit_tenant on admin_audit for select using (tenant_id = app_current_tenant_id());
create or replace function jsonb_shallow_diff(old jsonb, new jsonb)
returns jsonb language sql immutable as $$
  select coalesce(jsonb_object_agg(k, v), '{}'::jsonb) from (
    select k, new->k as v
    from (select jsonb_object_keys(coalesce(old,'{}'::jsonb) || coalesce(new,'{}'::jsonb)) as k) ks
    where coalesce(old->k,'null'::jsonb) is distinct from coalesce(new->k,'null'::jsonb)
  ) s;
$$;
create or replace function audit_entries() returns trigger language plpgsql as $$
begin
  if TG_OP='INSERT' then
    insert into admin_audit(tenant_id, actor, entity, entity_id, action, diff)
    values (new.tenant_id, coalesce(current_setting('request.jwt.claims',true),'system'), 'entries', new.id, 'insert', new.data);
  elsif TG_OP='UPDATE' then
    insert into admin_audit(tenant_id, actor, entity, entity_id, action, diff)
    values (new.tenant_id, coalesce(current_setting('request.jwt.claims',true),'system'), 'entries', new.id, 'update', jsonb_shallow_diff(old.data, new.data));
  elsif TG_OP='DELETE' then
    insert into admin_audit(tenant_id, actor, entity, entity_id, action, diff)
    values (old.tenant_id, coalesce(current_setting('request.jwt.claims',true),'system'), 'entries', old.id, 'delete', old.data);
  end if;
  return COALESCE(new, old);
end; $$;
drop trigger if exists trg_audit_entries on entries;
create trigger trg_audit_entries after insert or update or delete on entries for each row execute procedure audit_entries();
