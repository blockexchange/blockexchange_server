
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

-- add schemagroup_id column
alter table schema
  add column
  schemagroup_id serial;

-- add schemagroup for every user
insert into schemagroup(id, created, name)
  select id, created, name from public.user;

-- add user_schemagroup_permission
insert into user_schemagroup_permission(user_id, schemagroup_id, "create", "update", "delete", "manage")
  select id, id, true, true, true, true from public.user;

-- populate schemagroup_id
update schema set schemagroup_id = user_id;

-- add constraints to schemagroup
alter table schema add constraint schemagroup_id_fk
  foreign key (schemagroup_id) references schemagroup(id) on delete cascade;

-- drop user_id column
alter table schema drop user_id;
