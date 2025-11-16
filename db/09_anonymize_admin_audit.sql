-- Anonymize admin_audit actor fields and redact sensitive keys from diff JSONB.

create or replace function jsonb_strip_keys(in jb jsonb, in keys text[]) returns jsonb language plpgsql immutable as $$
declare
  k text;
  res jsonb := coalesce(jb, '{}'::jsonb);
begin
  if jb is null then return jb; end if;
  foreach k in array keys loop
    -- remove top-level key if exists
    res := res - k;
  end loop;
  return res;
end; $$;

create or replace function anonymize_admin_audit_for_user(p_tenant_id bigint, p_user_id text)
returns void language plpgsql security definer as $$
declare
  v_mask text := concat('anon:user:', encode(digest(p_user_id, 'sha256'), 'hex'));
  rec record;
begin
  for rec in select id, actor, diff from admin_audit where tenant_id = p_tenant_id loop
    -- if actor references the user_id, replace it
    if rec.actor is not null and position(p_user_id in rec.actor) > 0 then
      update admin_audit set actor = v_mask where id = rec.id;
    end if;
    -- redact common personal keys in diff (best-effort)
    if rec.diff is not null then
      update admin_audit set diff = jsonb_strip_keys(rec.diff, array['email','name','firstname','lastname']) where id = rec.id;
    end if;
  end loop;
end; $$;
