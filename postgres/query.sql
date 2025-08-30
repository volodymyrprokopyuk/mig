\set ON_ERROR_STOP on
\pset linestyle unicode
\pset columns 160
\pset expanded auto
\pset null ><

-- EU region
insert into eu.account(account_type, currency)
values ('payable', 'eur');
select acc.* from eu.account acc;

-- US region
insert into us.account(account_type, currency)
values ('payable', 'usd');
select acc.* from us.account acc;
