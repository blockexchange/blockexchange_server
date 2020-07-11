

create table schema_screenshot (
  id serial primary key not null,
  schema_id serial references schema(id) on delete cascade,
  title varchar(128) not null,
  data bytea not null
)
