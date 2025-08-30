do $$
declare
  a_schema constant text := '{{.Schema}}';
  a_version constant text := '00000';
begin

perform util.migration_check_applied(a_schema, a_version);

drop table {{.Schema}}.account;

drop type {{.Schema}}.currency_type;
drop type {{.Schema}}.account_type;

drop schema {{.Schema}};

perform util.migration_revert(a_schema, a_version);

end;
$$;
