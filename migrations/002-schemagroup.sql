

create table schemagroup(
  id serial primary key not null,
  created bigint not null,
  name varchar(64) unique
);

create table user_schemagroup_permission(
  user_id serial references public.user(id) on delete cascade,
  schemagroup_id serial references schemagroup(id) on delete cascade,
  "create" boolean,
  "update" boolean,
  "delete" boolean,
  "manage" boolean,
  primary key (user_id, schemagroup_id)
);
