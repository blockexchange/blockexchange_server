
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
