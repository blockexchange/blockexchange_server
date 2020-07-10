
create table public.user(
  id serial primary key not null,
  created bigint not null,
  name varchar not null,
  hash varchar not null,
  role varchar(16) not null,
  mail varchar
);

create unique index user_name on public.user(name);
create index user_created on public.user(created);


-- create temporary user with default password "temp"
insert into public.user(created, name, hash, role)
  values(
    extract(epoch from now()) * 1000,
    'temp',
    '$2a$10$g.6pRR93BwXfsMnPLWIKgOfIBDOcc48wJPCDtfNfzJbD/7zE2xgtm',
    'TEMP'
  );
