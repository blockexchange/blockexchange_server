
alter table schemapart add column mtime bigint not null default 0;
create index schemapart_id_mtime on schemapart(schema_id, mtime);
