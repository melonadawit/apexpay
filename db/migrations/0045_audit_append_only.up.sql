-- 0045_audit_append_only: make audit_logs truly append-only.
--
-- The SAD and marketing materials claim audit_logs is guarded by a trigger that
-- prevents UPDATE/DELETE, but no such trigger existed in the schema (security
-- audit F-1). Add it so historical audit rows cannot be mutated or deleted.

create or replace function prevent_audit_mutation() returns trigger as $$
begin
  raise exception 'audit_logs is append-only: UPDATE/DELETE not permitted';
end $$ language plpgsql;

drop trigger if exists audit_logs_no_mutation on audit_logs;
create trigger audit_logs_no_mutation
  before update or delete on audit_logs
  for each row execute function prevent_audit_mutation();
