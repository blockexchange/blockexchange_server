-- schemapart: serial-id -> uid
alter table schemapart add column uid uuid not null unique default gen_random_uuid();
alter table schemapart drop column id;

-- schemapart: order_id per schema_id
alter table schemapart add column order_id bigint;
update schemapart set order_id = offset_x + (offset_z * 2000) + (offset_y * 2000 * 2000);
alter table schemapart alter order_id set not null;
create index schemapart_order_id on schemapart(schema_id, order_id);

-- schemamod: serial-id -> uid
alter table schemamod add column uid uuid not null unique default gen_random_uuid();
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