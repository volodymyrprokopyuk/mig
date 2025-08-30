do $_$
declare
  a_schema constant text := '{{.Schema}}';
  a_version constant text := '00000';
  a_description constant text := 'Create {{.Schema}}.account';
begin

create schema {{.Schema}};

create type {{.Schema}}.account_type as enum (
  'receivable',
  'payable'
);

create type {{.Schema}}.currency_type as enum (
  {{ if (eq .Schema "us") }}
    'usd'
  {{ else }}
    'eur'
  {{ end }}
);

create table {{.Schema}}.account (
  id text not null default md5(random()::text),
  account_type {{.Schema}}.account_type not null,
  currency {{.Schema}}.currency_type not null,
  constraint pk_account primary key (id)
);

perform util.migration_apply(a_schema, a_version, a_description);

end;
$_$;
