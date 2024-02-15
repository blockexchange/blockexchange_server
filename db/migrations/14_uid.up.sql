-- schemapart: order_id per schema_id
alter table schemapart add column order_id bigint;
update schemapart set order_id = offset_x + (offset_z * 2000) + (offset_y * 2000 * 2000);
alter table schemapart alter order_id set not null;
alter table schemapart drop column id;

-- schemamod: serial-id -> uid
alter table schemamod drop column id;

-- schematag: serial-id -> uid
alter table schematag add column uid uuid not null unique default gen_random_uuid();
alter table schematag drop column id;

-- tag: serial-id -> uid
alter table tag add column uid uuid not null unique default gen_random_uuid();

-- relation: schematag -> tag
alter table schematag add column tag_uid uuid;
update schematag set tag_uid = (select t.uid from tag t where t.id = tag_id);
alter table schematag alter uid set not null;
alter table schematag drop column tag_id;
alter table tag drop column id;

-- access_token: serial-id -> uid
alter table access_token add column uid uuid not null unique default gen_random_uuid();
alter table access_token drop column id;

-- collection cleanup
drop table collection_schema;
drop table collection;

-- user: serial-id -> uid
alter table public.user add column uid uuid not null unique default gen_random_uuid();
alter table schema add column user_uid uuid default gen_random_uuid();
update schema set user_uid = (select u.uid from public.user u where u.id = user_id);
alter table schema alter user_uid set not null;
alter table schema drop column user_id;

alter table user_schema_star add column user_uid uuid default gen_random_uuid();
update user_schema_star set user_uid = (select u.uid from public.user u where u.id = user_id);
alter table user_schema_star alter user_uid set not null;
alter table user_schema_star drop column user_id;

alter table access_token add column user_uid uuid default gen_random_uuid();
update access_token set user_uid = (select u.uid from public.user u where u.id = user_id);
alter table access_token alter user_uid set not null;
alter table access_token drop column user_id;

alter table public.user drop column id;

-- schema_screenshot: serial-id -> uid
alter table schema_screenshot add column uid uuid not null unique default gen_random_uuid();
alter table schema_screenshot drop column id;

-- schema: serial-id -> uid
alter table schema add column uid uuid not null unique default gen_random_uuid();

alter table user_schema_star add column schema_uid uuid;
update user_schema_star set schema_uid = (select s.uid from schema s where s.id = schema_id);
alter table user_schema_star alter schema_uid set not null;
alter table user_schema_star drop column schema_id;

alter table schema_screenshot add column schema_uid uuid;
update schema_screenshot set schema_uid = (select s.uid from schema s where s.id = schema_id);
alter table schema_screenshot alter schema_uid set not null;
alter table schema_screenshot drop column schema_id;

alter table schemamod add column schema_uid uuid;
update schemamod set schema_uid = (select s.uid from schema s where s.id = schema_id);
alter table schemamod alter schema_uid set not null;
alter table schemamod drop column schema_id;
create unique index schemamod_schema_uid_mod_name on schemamod(schema_uid, mod_name);

alter table schemapart add column schema_uid uuid;
update schemapart set schema_uid = (select s.uid from schema s where s.id = schema_id);
alter table schemapart alter schema_uid set not null;
alter table schemapart drop column schema_id;
create unique index schemapart_offset_schema_uid on schemapart(offset_x, offset_y, offset_z, schema_uid);
create index schemapart_mtime_schema_uid on schemapart(mtime, schema_uid);
create index schemapart_order_id_schema_uid on schemapart(order_id, schema_uid);

alter table schematag add column schema_uid uuid;
update schematag set schema_uid = (select s.uid from schema s where s.id = schema_id);
alter table schematag alter schema_uid set not null;
alter table schematag drop column schema_id;

alter table schema drop column id;
