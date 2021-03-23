
-- USER Table

create table public.user(
  id serial primary key not null,
  created bigint not null,
  name varchar not null,
  hash varchar not null,
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
	token varchar not null,
	usecount int not null default 0
);

-- SCHEMA

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
	mtime bigint not null default 0,
  data bytea not null,
  metadata bytea not null
);

create unique index schemapart_coords on schemapart(schema_id, offset_x, offset_y, offset_z);

alter table schemapart
add constraint schemapart_unique_coords
unique using index schemapart_coords;

create index schemapart_id_mtime on schemapart(schema_id, mtime);

-- SCHEMAMOD

create table schemamod(
  id serial primary key not null,
  schema_id bigint not null references schema(id) on delete cascade,
  mod_name varchar(64) not null
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

-- COLLECTION

create table collection (
	id serial primary key not null,
  user_id serial references public.user(id) on delete cascade,
  name varchar not null
);

create table collection_schema (
	collection_id bigint not null references collection(id) on delete cascade,
	schema_id bigint not null references schema(id) on delete cascade,
	primary key (collection_id, schema_id)
);

-- TAG

create table tag(
	id serial primary key not null,
  name varchar(128) not null,
	description varchar not null
);

create table schematag(
  id serial primary key not null,
	tag_id bigint not null references tag(id) on delete cascade,
  schema_id bigint not null references schema(id) on delete cascade
);
