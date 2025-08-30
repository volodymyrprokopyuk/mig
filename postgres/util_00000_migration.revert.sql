do $$
declare
  a_schema constant text := 'util';
  a_version constant text := '00000';
begin

perform util.migration_check_applied(a_schema, a_version);

drop function util.migration_revert;
drop function util.migration_check_applied;
drop function util.migration_apply;
drop function util.migration_last_applied;

drop table util.migration;

drop schema util;

end;
$$;
