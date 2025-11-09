create or replace function set_updated_at() returns trigger language plpgsql as $$
begin new.updated_at:=now(); return new; end; $$;
create or replace function keep_version() returns trigger language plpgsql as $$
begin if TG_OP='UPDATE' then
  insert into entry_versions(entry_id,version,data,editor)
  values (old.id, coalesce((select max(version) from entry_versions where entry_id=old.id),0)+1, old.data, current_setting('request.jwt.claims',true));
end if; return new; end; $$;
create or replace function notify_content_changed() returns trigger language plpgsql as $$
begin
  if (TG_OP='UPDATE' and new.status='published' and old.status is distinct from 'published')
   or (TG_OP='INSERT' and new.status='published') then
    perform pg_notify('content_changed', json_build_object('tenant_id',new.tenant_id,'entry_id',new.id,'slug',new.slug)::text);
  end if; return new;
end; $$;
drop trigger if exists trg_entries_updated_at on entries;
create trigger trg_entries_updated_at before update on entries for each row execute procedure set_updated_at();
drop trigger if exists trg_entries_versioning on entries;
create trigger trg_entries_versioning before update on entries for each row execute procedure keep_version();
drop trigger if exists trg_entries_notify on entries;
create trigger trg_entries_notify after insert or update on entries for each row execute procedure notify_content_changed();
