do $_$
declare
  a_schema constant text := 'util';
  a_version constant text := '00000';
  a_description constant text := 'Create util.migration';
begin

set time zone 'UTC';

create schema util;

create table util.migration (
  schema text not null,
  version text not null,
  description text not null,
  applied timestamptz not null default now(),
  constraint pk_migration_version primary key (schema, version),
  constraint ck_migration_schema check (schema ~ '^\w{2,20}$'),
  constraint ck_migration_version check (version ~ '^\d{5}$')
);

create function util.migration_last_applied(a_schema text, a_count int)
returns setof util.migration stable
language plpgsql as $$
begin
  return query
  select mig.* from util.migration mig
  where mig.schema = a_schema
  order by mig.version desc limit a_count;
end;
$$;

create function util.migration_apply(
 a_schema text, a_version text, a_description text
)
returns void volatile
language plpgsql as $$
begin
  insert into util.migration(schema, version, description)
  values (a_schema, a_version, a_description);
end;
$$;

create function util.migration_check_applied(
  a_schema text, a_version text
)
returns void stable
language plpgsql as $$
begin
  perform 1 from util.migration mig
  where mig.schema = a_schema and mig.version = a_version;
  if not found then
    raise exception
      'version % for schema % is not applied', a_version, a_schema;
  end if;
end;
$$;

create function util.migration_revert(
  a_schema text, a_version text
)
returns void volatile
language plpgsql as $$
begin
  delete from util.migration
  where schema = a_schema and version = a_version;
end;
$$;

perform util.migration_apply(a_schema, a_version, a_description);

end;
$_$;
