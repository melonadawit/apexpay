-- 0045_audit_append_only.down.sql
drop trigger if exists audit_logs_no_mutation on audit_logs;
drop function if exists prevent_audit_mutation();
