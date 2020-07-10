
create table schema(
  id serial primary key not null,
  created bigint not null,
  schemagroup_id bigint not null references schemagroup(id) on delete cascade,
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

alter table schema add unique(schemagroup_id, name);

create index schema_created on schema(created);
