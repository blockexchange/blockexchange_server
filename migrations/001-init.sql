
create table public.user(
  id serial primary key not null,
  created bigint not null,
  name varchar not null,
  hash varchar not null,
  mail varchar
);

create unique index user_name on public.user(name);
create index user_created on public.user(created);


-- create temporary user with default password "temp"
insert into public.user(created, name, hash)
  values(
    extract(epoch from now()) * 1000,
    'temp',
    '$2a$10$g.6pRR93BwXfsMnPLWIKgOfIBDOcc48wJPCDtfNfzJbD/7zE2xgtm'
  );


create table schema(
  id serial primary key not null,
  created bigint not null,
  user_id bigint not null references public.user(id) on delete cascade,
  name varchar not null,
  description text not null,
  complete boolean not null,
  size_x smallint not null,
  size_y smallint not null,
  size_z smallint not null,
  part_length smallint not null,
  total_size int not null,
  total_parts int not null,
  search_tokens tsvector not null,
	downloads int not null default 0,
	license varchar not null default 'CC0',
	long_description text not null default ''
);

create index schema_created on schema(created);

create table schemapart(
  id serial primary key not null,
  schema_id bigint not null references schema(id) on delete cascade,
  offset_x smallint not null,
  offset_y smallint not null,
  offset_z smallint not null,
  data bytea not null,
  metadata bytea not null
);

create index schemapart_coords on schemapart(schema_id, offset_x, offset_y, offset_z);

create table schemamod(
  id serial primary key not null,
  schema_id bigint not null references schema(id) on delete cascade,
  mod_name varchar(64) not null,
  node_count int not null
);

create index schemamod_schema_id on schemamod(schema_id);
