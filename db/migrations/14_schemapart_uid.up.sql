alter table schemapart add column uid uuid not null unique default gen_random_uuid();
alter table schemapart drop column id;