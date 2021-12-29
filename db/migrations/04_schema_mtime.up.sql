alter table schema add column mtime bigint default 0 not null;

-- copy over create-date to mtime
update schema set mtime = created;