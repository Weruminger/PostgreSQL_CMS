-- GDPR helper: delete_user_gdpr
-- Attempts to remove or anonymize personal data related to a user for a tenant.
-- Returns a small JSON report { status, deleted_counts, errors }

create or replace function delete_user_gdpr(p_tenant_id bigint, p_user_id bigint)
returns jsonb language plpgsql security definer as $$
declare
  v_report jsonb := '{}';
  v_total_deleted int := 0;
  v_cnt int := 0;
  v_errors text[] := array[]::text[];
begin
  -- Run in a savepoint-friendly way
  perform set_config('app.tenant_id', p_tenant_id::text, true);

  -- Example: sessions table (if exists)
  begin
    execute format('delete from sessions where tenant_id = %s and user_id = %s', p_tenant_id, p_user_id);
    GET DIAGNOSTICS v_cnt = ROW_COUNT;
    v_total_deleted := v_total_deleted + v_cnt;
    v_report := jsonb_set(v_report, '{sessions}', to_jsonb(v_cnt), true);
  exception when others then
    v_errors := array_append(v_errors, 'sessions:' || sqlerrm);
  end;

  -- Example: media_assets owner cleanup (delete or set owner null)
  begin
    if exists (select 1 from information_schema.tables where table_name='media_assets') then
      execute format('update media_assets set meta = meta - ''owner'' where tenant_id = %s and (meta->>''owner_id'')::bigint = %s', p_tenant_id, p_user_id);
      GET DIAGNOSTICS v_cnt = ROW_COUNT;
      v_report := jsonb_set(v_report, '{media_assets_anonymized}', to_jsonb(v_cnt), true);
      v_total_deleted := v_total_deleted + v_cnt;
    end if;
  exception when others then
    v_errors := array_append(v_errors, 'media_assets:' || sqlerrm);
  end;

  -- Entries: remove personal fields from JSONB data (best-effort)
  begin
    if exists (select 1 from information_schema.tables where table_name='entries') then
      execute format(
        'update entries set data = (data - ''email'' - ''firstname'' - ''lastname'') where tenant_id = %s and (data->>''author_id'')::bigint = %s',
        p_tenant_id, p_user_id);
      GET DIAGNOSTICS v_cnt = ROW_COUNT;
      v_report := jsonb_set(v_report, '{entries_data_redacted}', to_jsonb(v_cnt), true);
      v_total_deleted := v_total_deleted + v_cnt;
    end if;
  exception when others then
    v_errors := array_append(v_errors, 'entries:' || sqlerrm);
  end;

  -- If a users table exists, delete the user row
  begin
    if exists (select 1 from information_schema.tables where table_name='users') then
      execute format('delete from users where tenant_id = %s and id = %s', p_tenant_id, p_user_id);
      GET DIAGNOSTICS v_cnt = ROW_COUNT;
      v_report := jsonb_set(v_report, '{users_deleted}', to_jsonb(v_cnt), true);
      v_total_deleted := v_total_deleted + v_cnt;
    end if;
  exception when others then
    v_errors := array_append(v_errors, 'users:' || sqlerrm);
  end;

  -- Anonymize related audit entries (call helper if exists)
  begin
    if exists (select 1 from pg_proc p join pg_namespace n on p.pronamespace = n.oid where p.proname = 'anonymize_admin_audit_for_user') then
      perform anonymize_admin_audit_for_user(p_tenant_id, p_user_id::text);
    end if;
  exception when others then
    v_errors := array_append(v_errors, 'anonymize_audit:' || sqlerrm);
  end;

  v_report := jsonb_set(v_report, '{status}', '"success"'::jsonb, true);
  v_report := jsonb_set(v_report, '{total_affected}', to_jsonb(v_total_deleted), true);
  if array_length(v_errors,1) is not null then
    v_report := jsonb_set(v_report, '{errors}', to_jsonb(v_errors), true);
  end if;
  return v_report;
end;
$$;
