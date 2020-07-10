
create table schemamod(
  id serial primary key not null,
  schema_id bigint not null references schema(id) on delete cascade,
  mod_name varchar(64) not null,
  node_count int not null
);

create index schemamod_schema_id on schemamod(schema_id);
