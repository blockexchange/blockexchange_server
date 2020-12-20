
-- USER Table

create table public.user(
  id serial primary key not null,
  created bigint not null,
  name varchar not null,
  hash varchar not null,
  role varchar(16) not null default 'MEMBER',
  type varchar(16) not null default 'LOCAL',
  external_id varchar(63),
  mail varchar
);

create unique index user_name on public.user(name);
create index user_created on public.user(created);

-- access_token

create table access_token(
	id serial primary key not null,
	user_id bigint not null references public.user(id) on delete cascade,
	created bigint not null,
	expires bigint not null,
	name varchar not null,
	token varchar not null
);

-- create temporary user with default password "temp"
insert into public.user(created, name, role, hash)
  values(
    extract(epoch from now()) * 1000,
    'temp',
    'UPLOAD_ONLY',
    '$2a$10$g.6pRR93BwXfsMnPLWIKgOfIBDOcc48wJPCDtfNfzJbD/7zE2xgtm'
  );

-- SCHEMA

create table schema(
  id serial primary key not null,
  created bigint not null,
  user_id bigint not null references public.user(id) on delete cascade,
  name varchar not null,
  description text not null,
  complete boolean not null,
  max_x smallint not null,
  max_y smallint not null,
  max_z smallint not null,
  part_length smallint not null,
  total_size int not null,
  total_parts int not null,
  search_tokens tsvector not null,
	downloads int not null default 0,
	license varchar not null default 'CC0'
);

alter table schema add unique(user_id, name);
create index schema_created on schema(created);

-- SCHEMAPART

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

-- SCHEMAMOD

create table schemamod(
  id serial primary key not null,
  schema_id bigint not null references schema(id) on delete cascade,
  mod_name varchar(64) not null,
  node_count int not null
);

create index schemamod_schema_id on schemamod(schema_id);

-- USER_SCHEMA_STAR

create table user_schema_star (
  user_id serial references public.user(id) on delete cascade,
  schema_id serial references schema(id) on delete cascade,
  primary key (user_id, schema_id)
);

-- SCHEMA_SCREENSHOT

create table schema_screenshot (
  id serial primary key not null,
  schema_id serial references schema(id) on delete cascade,
	type varchar(64) not null,
  title varchar(128) not null,
  data bytea not null
);
