
-- add role column
alter table public.user
  add column
  role varchar(16);

-- set default value
update public.user
  set role = 'MEMBER';

-- set temp user role to "temp"
update public.user
  set role = 'TEMP'
  where name = 'temp';

-- non-null constraint
alter table public.user
  alter column
  role set not null;
