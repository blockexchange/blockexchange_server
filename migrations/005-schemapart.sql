
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
