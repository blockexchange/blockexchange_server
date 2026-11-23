-- mod/game dependencies of a schema
create table schema_mod_dependency(
    uid uuid primary key not null default gen_random_uuid(),
    schema_uid uuid not null references schema(uid) on delete cascade,
    mod_type varchar(16) not null check (mod_type in ('mod', 'game')),
    git_url varchar(256) not null,
    version varchar(64) not null default 'latest'
);
create index schema_mod_dependency_schema_uid on schema_mod_dependency(schema_uid);
