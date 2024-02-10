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
